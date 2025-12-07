import CryptoJS from "crypto-js";

const AES_SECRET = "12345678901234567890123456789012";

// --------------------------------------------
// Recreate OpenSSL EVP_BytesToKey (SHA256)
// --------------------------------------------
function evpBytesToKey(password, salt) {
  let data = CryptoJS.enc.Utf8.parse("");
  let key = CryptoJS.enc.Utf8.parse("");
  
  let prev = CryptoJS.enc.Utf8.parse("");

  while (key.sigBytes < 48) {
    const md = CryptoJS.algo.SHA256.create();
    md.update(prev);
    md.update(password);
    md.update(salt);
    prev = md.finalize();
    key = key.concat(prev);
  }

  const keyBytes = CryptoJS.lib.WordArray.create(key.words.slice(0, 8), 32);  // 32 bytes
  const ivBytes  = CryptoJS.lib.WordArray.create(key.words.slice(8, 12), 16); // 16 bytes

  return { key: keyBytes, iv: ivBytes };
}

// --------------------------------------------
// Decrypt PKCS8 Private Key from Backend (AES-256-CBC + Salted__)
// --------------------------------------------
export function decryptPrivateKey(b64) {
  if (!b64) throw new Error("Empty encrypted key");

  console.log("Encrypted key from Firestore:", b64);

  const raw = CryptoJS.enc.Base64.parse(b64);
  const rawBytes = raw.toString(CryptoJS.enc.Latin1);

  if (rawBytes.slice(0, 8) !== "Salted__") {
    throw new Error("Invalid OpenSSL data (missing Salted__)");
  }

  const salt = CryptoJS.enc.Latin1.parse(rawBytes.slice(8, 16));
  const ciphertext = CryptoJS.enc.Latin1.parse(rawBytes.slice(16));

  // Derive key + iv exactly as Go does
  const { key, iv } = evpBytesToKey(
    CryptoJS.enc.Utf8.parse(AES_SECRET),
    salt
  );

  // AES-256-CBC decrypt
  const decrypted = CryptoJS.AES.decrypt(
    { ciphertext },
    key,
    {
      iv,
      mode: CryptoJS.mode.CBC,
      padding: CryptoJS.pad.Pkcs7
    }
  );

  const plaintext = decrypted.toString(CryptoJS.enc.Utf8);
  if (!plaintext) {
    throw new Error("AES decryption failed: wrong password or corrupted data");
  }

  return plaintext;  // This is still BASE64 PKCS8 private key
}
