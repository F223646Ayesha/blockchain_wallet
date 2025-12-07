export async function importPrivateKeyEcdsaP256(pkcs8Base64) {

  console.log("==================================");
  console.log("🔵 importPrivateKeyEcdsaP256 CALLED");
  console.log("🔵 PKCS8 (base64):", pkcs8Base64);
  console.log("==================================");

  // Convert base64 → ArrayBuffer
  const binary = atob(pkcs8Base64);
  const bytes = new Uint8Array(binary.length);
  for (let i = 0; i < binary.length; i++) bytes[i] = binary.charCodeAt(i);

  return crypto.subtle.importKey(
    "pkcs8",
    bytes.buffer,
    { name: "ECDSA", namedCurve: "P-256" },
    true,
    ["sign"]
  );
}
