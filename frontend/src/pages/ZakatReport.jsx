import { useEffect, useState } from "react";
import { apiGet } from "../api/api";

export default function ZakatReport() {
  const [history, setHistory] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  const walletId = localStorage.getItem("walletId");

  useEffect(() => {
    async function load() {
      try {
        const res = await apiGet(`/transaction/history/${walletId}`);

        if (!res.success) {
          setError("Failed to load zakat report");
        } else {
          const zakatTx = res.data.filter(
            (tx) =>
              tx.type === "sent" &&
              (tx.note === "Monthly zakat deduction" ||
                tx.note?.toLowerCase().includes("zakat"))
          );
          setHistory(zakatTx);
        }
      } catch (err) {
        setError("Backend not reachable");
      }
      setLoading(false);
    }

    load();
  }, []);

  return (
    <div className="max-w-4xl mx-auto py-10 animate-fadeIn">
      {/* Title */}
      <h1
        className="
          text-4xl font-bold mb-8 text-transparent bg-clip-text 
          bg-gradient-to-r from-blue-400 to-teal-300
        "
      >
        Zakat Deduction Report
      </h1>

      {/* Loading */}
      {loading && (
        <div
          className="
            p-4 rounded-xl border border-white/10 bg-white/5 
            text-gray-300 backdrop-blur-xl
          "
        >
          Loading zakat history...
        </div>
      )}

      {/* Error */}
      {error && (
        <div
          className="
            mt-4 p-4 rounded-xl bg-red-900/30 border border-red-500/30 
            text-red-300 backdrop-blur-xl
          "
        >
          {error}
        </div>
      )}

      {/* No Data */}
      {!loading && history.length === 0 && !error && (
        <div
          className="
            p-4 rounded-xl border border-white/10 bg-white/5 
            text-gray-400 text-sm backdrop-blur-xl
          "
        >
          No zakat deductions found.
        </div>
      )}

      {/* Zakat Records */}
      <div className="space-y-4 mt-4">
        {history.map((tx, i) => (
          <div
            key={i}
            className="
              p-5 rounded-2xl bg-gradient-to-br from-gray-900/60 to-gray-800/40 
              border border-white/10 shadow-xl backdrop-blur-xl
              flex items-center justify-between
              hover:shadow-blue-500/20 transition transform hover:-translate-y-1
            "
          >
            {/* Left Section */}
            <div>
              <p className="text-lg font-semibold text-blue-300">
                {tx.amount} coins
              </p>
              <p className="text-gray-400 text-sm">{tx.note || "Zakat"}</p>
            </div>

            {/* Timestamp */}
            <p className="text-gray-500 text-xs">
              {new Date(tx.timestamp * 1000).toLocaleString()}
            </p>
          </div>
        ))}
      </div>
    </div>
  );
}
