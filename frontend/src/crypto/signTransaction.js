import { importPrivateKeyEcdsaP256 } from "./keyUtils";

// ---------------------------------------------------
// Convert bytes → hex
// ---------------------------------------------------
function bytesToHex(bytes) {
  return [...bytes].map((b) => b.toString(16).padStart(2, "0")).join("");
}

// ---------------------------------------------------
// Convert raw WebCrypto signature (R||S) → ASN.1 DER
// Go expects ASN1, not raw 64-byte signature
// ---------------------------------------------------
function rawToAsn1(rawSig) {
  console.log("🟡 rawToAsn1() called");

  const r = rawSig.slice(0, 32);
  const s = rawSig.slice(32);

  console.log("🟡 R part:", r);
  console.log("🟡 S part:", s);

  function trimLeadingZeros(bytes) {
    let i = 0;
    while (i < bytes.length - 1 && bytes[i] === 0) i++;
    return bytes.slice(i);
  }

  let rTrim = trimLeadingZeros(r);
  let sTrim = trimLeadingZeros(s);

  console.log("🟡 R trimmed:", rTrim);
  console.log("🟡 S trimmed:", sTrim);

  // If MSB = 1, prepend 0x00 to avoid negative INTEGER
  function encodeInt(bytes) {
    if (bytes[0] & 0x80) {
      bytes = Uint8Array.from([0, ...bytes]);
    }
    return Uint8Array.from([0x02, bytes.length, ...bytes]);
  }

  const rEnc = encodeInt(rTrim);
  const sEnc = encodeInt(sTrim);

  console.log("🟡 ASN.1 R encoded:", rEnc);
  console.log("🟡 ASN.1 S encoded:", sEnc);

  const seq = Uint8Array.from([
    0x30,
    rEnc.length + sEnc.length,
    ...rEnc,
    ...sEnc,
  ]);

  console.log("🟡 FINAL ASN.1 SIGNATURE (bytes):", seq);
  return seq;
}

// ---------------------------------------------------
// WebCrypto signer
// ---------------------------------------------------
export async function signData(privateKeyObj, message) {
  if (!privateKeyObj) throw new Error("Private key not loaded");

  console.log("🔵 signData() called");
  console.log("🔵 MESSAGE TO SIGN:", message);

  const data = new TextEncoder().encode(message);
  console.log("🔵 MESSAGE BYTES:", data);

  // RAW WebCrypto signature (R||S)
  const raw = new Uint8Array(
    await crypto.subtle.sign(
      { name: "ECDSA", hash: "SHA-256" },
      privateKeyObj,
      data
    )
  );

  console.log("🔵 RAW SIGNATURE (R||S):", raw);

  // Convert → ASN.1 DER for Go
  const asn1 = rawToAsn1(raw);

  console.log("🔵 ASN.1 DER SIGNATURE:", asn1);

  const hex = bytesToHex(asn1);

  console.log("🔵 FINAL SIGNATURE (HEX SENT TO BACKEND):", hex);

  return hex;
}

// ---------------------------------------------------
// MUST MATCH BACKEND buildPayload EXACTLY
// ---------------------------------------------------
function buildPayload(tx) {
  const payload =
    tx.sender +
    "|" +
    tx.receiver +
    "|" +
    tx.amount.toString() +
    "|" +
    tx.timestamp.toString() +
    "|" +
    tx.note;

  console.log("🟣 buildPayload() OUTPUT:", payload);
  return payload;
}

// ---------------------------------------------------
// Main entry — used by SendMoney.jsx
// ---------------------------------------------------
export async function generateTransactionSignature(tx) {
  console.log("========================================");
  console.log("🔶 generateTransactionSignature() CALLED");

  const privateKeyHex = localStorage.getItem("privateKeyHex");
  console.log("🔶 privateKeyHex (PKCS8 base64):", privateKeyHex);

  if (!privateKeyHex) throw new Error("Missing decrypted private key");

  // Import key
  const privateKeyObj = await importPrivateKeyEcdsaP256(privateKeyHex);
  console.log("🔶 WebCrypto PrivateKey object imported OK");

  // Build payload
  const payload = buildPayload(tx);
  console.log("🔶 FINAL PAYLOAD for SIGNING:", payload);

  // Sign (ASN.1 encoded)
  const signature = await signData(privateKeyObj, payload);

  console.log("🔶 SIGNATURE RETURNED:", signature);
  console.log("========================================");

  return signature;
}
