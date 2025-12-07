import { useEffect, useState } from "react";
import { apiGet, apiPost } from "../api/api";
import { generateTransactionSignature } from "../crypto/signTransaction";

export default function SendMoney() {
  const [wallet, setWallet] = useState(null);
  const [beneficiaries, setBeneficiaries] = useState([]);

  const [receiver, setReceiver] = useState("");
  const [amount, setAmount] = useState("");
  const [note, setNote] = useState("");

  const [loading, setLoading] = useState(false);
  const [result, setResult] = useState(null);
  const [error, setError] = useState("");

  // ===========================================================
  // LOAD WALLET + BENEFICIARIES
  // ===========================================================
  useEffect(() => {
    async function loadWallet() {
      try {
        const walletId = localStorage.getItem("walletId");
        if (!walletId) return;

        const res = await apiGet(`/wallet/profile/${walletId}`);

        if (!res.success) {
          setError("Failed to load wallet");
          return;
        }

        const w = res.data;
        setWallet(w);
        setBeneficiaries(w.beneficiaries || []);

        // store signing keys
        if (w.public_key) localStorage.setItem("publicKey", w.public_key);
        if (w.private_key_enc) localStorage.setItem("privateKeyEnc", w.private_key_enc);

      } catch (err) {
        setError("Error loading wallet");
      }
    }

    loadWallet();
  }, []);

  // ===========================================================
  // SEND MONEY
  // ===========================================================
  async function handleSend(e) {
    e.preventDefault();
    setLoading(true);
    setError("");
    setResult(null);

    try {
      const sender = wallet.wallet_id;
      const timestamp = Math.floor(Date.now() / 1000);

      const tx = {
        sender,
        receiver,
        amount: parseFloat(amount),
        note,
        timestamp,
      };

      const signature = await generateTransactionSignature(tx);
      tx.signature = signature;

      const res = await apiPost("/transaction/send", tx);

      if (!res.success) {
        setError(res.message || "Transaction failed");
      } else {
        setResult(res);
        setAmount("");
        setNote("");
        setReceiver("");
      }

    } catch (err) {
      setError("Transaction failed: " + (err.message || err));
    }

    setLoading(false);
  }

  // ===========================================================
  // UI
  // ===========================================================
  return (
    <div className="max-w-4xl mx-auto py-10 animate-fadeIn">

      <h1 className="text-4xl font-bold text-blue-400 mb-8 drop-shadow">
        Send Money
      </h1>

      {!wallet ? (
        <div className="p-4 bg-yellow-900/40 text-yellow-300 border border-yellow-700 rounded">
          Loading wallet...
        </div>
      ) : (
        <form
          onSubmit={handleSend}
          className="
            bg-white/5 backdrop-blur-xl border border-white/10
            p-8 rounded-2xl shadow-lg space-y-6
          "
        >

          {/* Receiver */}
          <div>
            <label className="text-gray-300 text-sm">Receiver Wallet</label>
            <select
              className="
                w-full mt-1 p-3 rounded-lg bg-black/40 border border-gray-700 
                text-gray-200 focus:ring-2 focus:ring-blue-500 outline-none
              "
              value={receiver}
              onChange={(e) => setReceiver(e.target.value)}
              required
            >
              <option value="">Select Beneficiary</option>
              {beneficiaries.map((b) => (
                <option key={b} value={b}>{b}</option>
              ))}
            </select>
          </div>

          {/* Amount */}
          <div>
            <label className="text-gray-300 text-sm">Amount</label>
            <input
              type="number"
              className="
                w-full mt-1 p-3 rounded-lg bg-black/40 border border-gray-700 
                text-gray-200 focus:ring-2 focus:ring-green-500 outline-none
              "
              value={amount}
              onChange={(e) => setAmount(e.target.value)}
              required
            />
          </div>

          {/* Note */}
          <div>
            <label className="text-gray-300 text-sm">Note</label>
            <input
              type="text"
              className="
                w-full mt-1 p-3 rounded-lg bg-black/40 border border-gray-700 
                text-gray-200 focus:ring-2 focus:ring-purple-500 outline-none
              "
              value={note}
              onChange={(e) => setNote(e.target.value)}
              placeholder="Optional message"
            />
          </div>

          {/* SEND BTN */}
          <button
            disabled={loading}
            className="
              w-full py-3 rounded-xl text-white font-semibold
              bg-gradient-to-r from-blue-600 to-purple-600
              hover:from-blue-500 hover:to-purple-500
              shadow-lg hover:shadow-purple-500/30 transition
            "
          >
            {loading ? "Processing..." : "Send Money"}
          </button>
        </form>
      )}

      {/* ERROR */}
      {error && (
        <div className="mt-4 p-4 bg-red-900/40 text-red-300 border border-red-700 rounded">
          {error}
        </div>
      )}

      {/* SUCCESS RESULT */}
      {result && (
        <pre className="
          mt-6 p-4 bg-green-900/40 text-green-300 border border-green-700
          rounded overflow-auto text-sm whitespace-pre-wrap
        ">
          {JSON.stringify(result, null, 2)}
        </pre>
      )}
    </div>
  );
}
