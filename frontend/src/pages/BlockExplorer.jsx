import { useEffect, useState } from "react";
import { apiGet } from "../api/api";

export default function BlockExplorer() {
  const [blocks, setBlocks] = useState([]);
  const [loading, setLoading] = useState(true);
  const [validationMsg, setValidationMsg] = useState("");
  const [validationError, setValidationError] = useState("");
  const [expanded, setExpanded] = useState({}); // toggle block open/close

  useEffect(() => {
    async function load() {
      try {
        const res = await apiGet("/blockchain");
        if (res.success) setBlocks(res.data || []);
      } catch (err) {
        console.error(err);
      }
      setLoading(false);
    }

    load();
  }, []);

  // ----------------- VALIDATE BLOCKCHAIN -----------------
  async function validateChain() {
    setValidationMsg("");
    setValidationError("");

    try {
      const res = await apiGet("/blockchain/validate");

      if (res.success) {
        setValidationMsg(res.message || "Blockchain is valid!");
      } else {
        setValidationError(res.message || "Blockchain invalid!");
      }
    } catch (err) {
      setValidationError("Error validating blockchain");
    }
  }

  function toggleExpand(index) {
    setExpanded((prev) => ({ ...prev, [index]: !prev[index] }));
  }

  if (loading)
    return (
      <p className="text-gray-300 text-center mt-10 animate-pulse">
        Loading blockchain...
      </p>
    );

  return (
    <div className="max-w-6xl mx-auto p-8 animate-fadeIn">

      {/* Title + Validate Button */}
      <div className="flex items-center justify-between mb-8">
        <h1 className="text-4xl font-bold text-blue-400 drop-shadow">
          Blockchain Explorer
        </h1>

        <button
          onClick={validateChain}
          className="
            px-5 py-2 rounded-xl text-white font-semibold
            bg-gradient-to-r from-green-600 to-emerald-700
            hover:from-green-500 hover:to-emerald-600
            shadow-lg hover:shadow-emerald-500/30
            transition
          "
        >
          Validate Blockchain
        </button>
      </div>

      {/* Validation result */}
      {validationMsg && (
        <div className="p-4 mb-4 rounded-xl text-green-300 bg-green-900/30 border border-green-700">
          ✅ {validationMsg}
        </div>
      )}

      {validationError && (
        <div className="p-4 mb-4 rounded-xl text-red-300 bg-red-900/30 border border-red-700">
          ❌ {validationError}
        </div>
      )}

      {/* ---------------- BLOCK LIST ---------------- */}
      <div className="space-y-6">
        {blocks.map((b, i) => (
          <div
            key={i}
            className="
              bg-white/5 backdrop-blur-xl border border-white/10
              rounded-2xl p-6 shadow-lg hover:shadow-blue-500/20
              transition cursor-pointer
            "
          >
            {/* ---------------- HEADER ---------------- */}
            <div
              onClick={() => toggleExpand(i)}
              className="flex justify-between items-center"
            >
              <h2 className="text-xl font-semibold text-blue-300">
                Block #{b.index}
              </h2>

              <button
                className="text-gray-300 text-sm hover:text-white transition"
              >
                {expanded[i] ? "Hide ▲" : "View ▼"}
              </button>
            </div>

            {/* ---------------- COLLAPSIBLE CONTENT ---------------- */}
            {expanded[i] && (
              <div className="mt-4 space-y-3">

                {/* Hash */}
                <div>
                  <p className="text-gray-400 text-sm">Hash</p>
                  <div className="w-full p-2 rounded bg-black/40 border border-gray-700 text-gray-200 font-mono break-all text-xs">
                    {b.hash}
                  </div>
                </div>

                {/* Previous Hash */}
                <div>
                  <p className="text-gray-400 text-sm">Previous Hash</p>
                  <div className="w-full p-2 rounded bg-black/40 border border-gray-700 text-gray-200 font-mono break-all text-xs">
                    {b.previous_hash}
                  </div>
                </div>

                {/* Merkle Root */}
                <div>
                  <p className="text-gray-400 text-sm">Merkle Root</p>
                  <div className="w-full p-2 rounded bg-black/40 border border-gray-700 text-gray-200 font-mono break-all text-xs">
                    {b.merkle_root}
                  </div>
                </div>

                {/* Nonce */}
                <div>
                  <p className="text-gray-400 text-sm">Nonce</p>
                  <p className="text-white">{b.nonce}</p>
                </div>

                {/* Transactions */}
                <div className="mt-4">
                  <p className="text-gray-400 text-sm mb-1">Transactions</p>
                  <pre className="
                    bg-black/30 border border-gray-700 p-4 rounded-xl
                    text-gray-200 text-xs overflow-auto
                    max-h-64
                  ">
                    {JSON.stringify(b.transactions, null, 2)}
                  </pre>
                </div>

              </div>
            )}
          </div>
        ))}
      </div>
    </div>
  );
}
