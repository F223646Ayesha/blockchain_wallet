import { useEffect, useState } from "react";
import { apiGet, apiPost } from "../api/api";

export default function Dashboard() {
  const [user, setUser] = useState(null);
  const [walletId, setWalletId] = useState(localStorage.getItem("walletId"));
  const [wallet, setWallet] = useState(null);

  const [loading, setLoading] = useState(false);
  const [creating, setCreating] = useState(false);
  const [mining, setMining] = useState(false);
  const [error, setError] = useState("");

  useEffect(() => {
    const stored = localStorage.getItem("user");
    if (stored) {
      const parsed = JSON.parse(stored);
      setUser(parsed);
      if (parsed.wallet_id) {
        setWalletId(parsed.wallet_id);
        localStorage.setItem("walletId", parsed.wallet_id);
      }
    }
  }, []);

  const loadWallet = async () => {
    if (!walletId) {
      setError("No wallet linked to this user.");
      return;
    }

    try {
      setLoading(true);
      setError("");
      const res = await apiGet(`/wallet/${walletId}`);
      if (!res.success) setError(res.message || "Invalid wallet ID");
      else setWallet(res.data);
    } catch {
      setError("Backend unreachable. Is Go server running?");
    }
    setLoading(false);
  };

  useEffect(() => {
    if (walletId) loadWallet();
  }, [walletId]);

  const handleCreateNew = async () => {
    try {
      setCreating(true);
      const res = await apiPost("/wallet/create", {});
      if (!res.success) return setError(res.message || "Wallet creation failed");

      const data = res.data;
      localStorage.setItem("walletId", data.wallet_id);
      localStorage.setItem("publicKey", data.public_key);
      localStorage.setItem("privateKeyEnc", data.private_key_encrypted);

      setWalletId(data.wallet_id);
      setWallet(data);
    } catch {
      setError("Error creating wallet");
    }
    setCreating(false);
  };

  const mineBlock = async () => {
    setError("");
    try {
      setMining(true);
      const res = await apiPost(`/mine?miner=${walletId}`);
      if (!res.success) return setError(res.message || "Mining failed");

      const block = res.data;
      alert(`⛏️ Block Mined!\nHash: ${block.hash}`);
      loadWallet();
    } catch {
      setError("Mining failed. Backend unreachable.");
    }
    setMining(false);
  };

  return (
    <div className="min-h-screen px-6 py-8 bg-gradient-to-b from-[#0a0f1f] to-black text-gray-100">
      
      <h1 className="text-4xl font-bold text-blue-400 drop-shadow-md mb-8 animate-fadeIn">
        Dashboard
      </h1>

      <div className="bg-white/5 backdrop-blur-xl p-6 rounded-2xl border border-white/10 shadow-2xl transition hover:shadow-blue-500/20 animate-slideUp">

        <h2 className="text-2xl font-semibold text-blue-300 mb-4">Your Wallet</h2>

        {!walletId ? (
          <p className="text-gray-400">No wallet found for this user.</p>
        ) : (
          <div className="space-y-4">
            
            <div>
              <p className="text-gray-400 text-sm">Wallet ID</p>
              <p className="font-mono break-words bg-black/40 p-3 rounded border border-gray-700">
                {walletId}
              </p>
            </div>

            {wallet && (
              <>
                <div>
                  <p className="text-gray-400 text-sm">Public Key</p>
                  <p className="font-mono break-words bg-black/40 p-3 rounded border border-gray-700 text-xs">
                    {wallet.public_key}
                  </p>
                </div>

                <div>
                  <p className="text-gray-400 text-sm">Encrypted Private Key</p>
                  <p className="font-mono break-words bg-black/40 p-3 rounded border border-gray-700 text-xs">
                    {wallet.private_key_enc}
                  </p>
                </div>
              </>
            )}
          </div>
        )}

        {error && (
          <div className="mt-4 p-3 bg-red-900/40 border border-red-600/40 text-red-300 rounded">
            {error}
          </div>
        )}

        <div className="flex gap-4 mt-6">
          <button
            onClick={loadWallet}
            disabled={loading}
            className="px-4 py-2 rounded-lg bg-blue-600 hover:bg-blue-500 transition shadow hover:shadow-blue-400/30"
          >
            {loading ? "Refreshing..." : "Refresh Wallet"}
          </button>

          <button
            onClick={handleCreateNew}
            disabled={creating}
            className="px-4 py-2 rounded-lg bg-green-600 hover:bg-green-500 transition shadow hover:shadow-green-400/30"
          >
            {creating ? "Creating..." : "Create New Wallet"}
          </button>

          <button
            onClick={mineBlock}
            disabled={mining}
            className="px-4 py-2 rounded-lg bg-purple-600 hover:bg-purple-500 transition shadow hover:shadow-purple-400/30"
          >
            {mining ? "Mining..." : "⛏️ Mine Block"}
          </button>
        </div>
      </div>

      {wallet && (
        <div className="mt-10 bg-black/40 backdrop-blur-xl p-6 rounded-xl border border-white/10 shadow-xl animate-fadeIn">
          <h2 className="text-xl font-semibold mb-3 text-blue-300">Wallet JSON</h2>
          <pre className="whitespace-pre-wrap text-xs bg-black/60 p-4 rounded-lg border border-gray-700 break-words">
            {JSON.stringify(wallet, null, 2)}
          </pre>
        </div>
      )}
    </div>
  );
}
