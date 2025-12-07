import { useEffect, useState } from "react";
import { apiGet } from "../api/api";

export default function TransactionHistory() {
  const [history, setHistory] = useState([]);
  const [loading, setLoading] = useState(true);

  const walletId = localStorage.getItem("walletId");

  useEffect(() => {
    async function load() {
      try {
        const res = await apiGet(`/transaction/history/${walletId}`);
        if (res.success) setHistory(res.data || []);
      } catch (err) {
        console.error("TX History Error:", err);
      }
      setLoading(false);
    }

    load();
  }, [walletId]);

  if (loading)
    return (
      <p className="text-center text-gray-400 mt-10 animate-pulse">
        Loading transaction history...
      </p>
    );

  if (!history.length)
    return (
      <p className="text-center text-gray-500 mt-10">
        No transactions found.
      </p>
    );

  return (
    <div className="max-w-5xl mx-auto p-8 animate-fadeIn space-y-6">
      {/* Page Title */}
      <h1
        className="text-4xl font-bold mb-6 text-transparent bg-clip-text 
        bg-gradient-to-r from-purple-400 via-pink-400 to-blue-400"
      >
        Transaction History
      </h1>

      {/* Transaction Cards */}
      {history.map((tx, i) => {
        const type = tx.type || "unknown"; // SAFE fallback

        const isSent = type === "sent";
        const isReceived = type === "received";
        const isZakat = tx.note?.toLowerCase().includes("zakat");

        let amountColor = "text-gray-300";
        if (isSent) amountColor = "text-red-400";
        if (isReceived) amountColor = "text-green-400";
        if (isZakat) amountColor = "text-blue-400";

        return (
          <div
            key={i}
            className="
              p-6 rounded-2xl border border-white/10 bg-white/5 
              backdrop-blur-xl shadow-xl hover:shadow-purple-500/20
              transition transform hover:-translate-y-1
            "
          >
            {/* First row: Amount + Type + Timestamp */}
            <div className="flex items-center justify-between mb-4">
              <p className={`text-2xl font-bold ${amountColor}`}>
                {tx.amount ?? 0} coins
              </p>

              <span
                className={`
                  px-3 py-1 rounded-full text-xs font-semibold
                  ${
                    isZakat &&
                    "bg-blue-500/20 text-blue-300 border border-blue-500/30"
                  }
                  ${
                    isSent &&
                    "bg-red-500/20 text-red-300 border border-red-500/30"
                  }
                  ${
                    isReceived &&
                    "bg-green-500/20 text-green-300 border border-green-500/30"
                  }
                  ${
                    !isSent && !isReceived && !isZakat &&
                    "bg-gray-500/20 text-gray-300 border border-gray-500/30"
                  }
                `}
              >
                {isZakat ? "Zakat" : type.toUpperCase()}
              </span>

              <p className="text-gray-400 text-sm">
                {tx.timestamp
                  ? new Date(tx.timestamp * 1000).toLocaleString()
                  : "N/A"}
              </p>
            </div>

            {/* Sender + Receiver */}
            <div className="grid md:grid-cols-2 gap-4 text-sm">
              <div>
                <p className="text-gray-400">Sender</p>
                <p className="font-mono text-gray-200 break-all bg-black/20 p-2 rounded-lg border border-white/10">
                  {tx.sender || "N/A"}
                </p>
              </div>

              <div>
                <p className="text-gray-400">Receiver</p>
                <p className="font-mono text-gray-200 break-all bg-black/20 p-2 rounded-lg border border-white/10">
                  {tx.receiver || "N/A"}
                </p>
              </div>
            </div>

            {/* Optional Note */}
            {tx.note && (
              <p className="mt-4 text-gray-300">
                <span className="font-semibold text-gray-400">Note: </span>
                {tx.note}
              </p>
            )}
          </div>
        );
      })}
    </div>
  );
}
