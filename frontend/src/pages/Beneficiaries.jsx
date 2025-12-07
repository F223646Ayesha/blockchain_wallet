import { useEffect, useState } from "react";
import { apiGet, apiPost } from "../api/api";

export default function Beneficiaries() {
  const [wallet, setWallet] = useState(null);
  const [beneficiaries, setBeneficiaries] = useState([]);
  const [newBeneficiary, setNewBeneficiary] = useState("");

  const [loading, setLoading] = useState(false);
  const [pageLoading, setPageLoading] = useState(true);
  const [error, setError] = useState("");
  const [success, setSuccess] = useState("");

  // ------------------------------------------------------
  // LOAD WALLET + BENEFICIARIES
  // ------------------------------------------------------
  useEffect(() => {
    async function loadData() {
      try {
        const walletId = localStorage.getItem("walletId");
        if (!walletId) {
          setError("No wallet linked to this account.");
          setPageLoading(false);
          return;
        }

        const res = await apiGet(`/wallet/${walletId}`);

        if (!res.success) {
          setError(res.message || "Failed to load wallet");
          setPageLoading(false);
          return;
        }

        setWallet(res.data);
        setBeneficiaries(res.data.beneficiaries || []);
      } catch (err) {
        console.error(err);
        setError("Failed to load wallet");
      }
      setPageLoading(false);
    }
    loadData();
  }, []);

  // ------------------------------------------------------
  // ADD BENEFICIARY
  // ------------------------------------------------------
  async function addBeneficiary() {
    setError("");
    setSuccess("");

    if (!wallet) {
      setError("Wallet not loaded yet.");
      return;
    }

    if (!newBeneficiary.trim()) {
      setError("Enter a valid wallet ID");
      return;
    }

    try {
      setLoading(true);

      const body = {
        wallet_id: wallet.wallet_id,
        beneficiary_id: newBeneficiary.trim(),
      };

      const res = await apiPost("/wallet/beneficiary/add", body);

      if (!res.success) {
        setError(res.message || "Failed to add beneficiary");
      } else {
        setSuccess("Beneficiary added successfully");
        setBeneficiaries([...beneficiaries, newBeneficiary.trim()]);
        setNewBeneficiary("");
      }
    } catch (err) {
      console.error(err);
      setError("Something went wrong while adding beneficiary");
    }

    setLoading(false);
  }

  // ------------------------------------------------------
  // REMOVE BENEFICIARY
  // ------------------------------------------------------
  async function removeBeneficiary(b) {
    setError("");
    setSuccess("");

    if (!wallet) {
      setError("Wallet not loaded yet.");
      return;
    }

    try {
      setLoading(true);

      const body = {
        wallet_id: wallet.wallet_id,
        beneficiary_id: b,
      };

      const res = await apiPost("/wallet/beneficiary/remove", body);

      if (!res.success) {
        setError(res.message || "Failed to remove beneficiary");
      } else {
        setSuccess("Beneficiary removed");
        setBeneficiaries(beneficiaries.filter((x) => x !== b));
      }
    } catch (err) {
      console.error(err);
      setError("Something went wrong while removing beneficiary");
    }

    setLoading(false);
  }

  // ------------------------------------------------------
  // UI
  // ------------------------------------------------------
  if (pageLoading) {
    return (
      <div className="min-h-[60vh] flex items-center justify-center">
        <p className="text-gray-400 animate-pulse">Loading beneficiaries...</p>
      </div>
    );
  }

  return (
    <div className="max-w-4xl mx-auto py-10 animate-fadeIn">
      {/* Header */}
      <div className="flex flex-col md:flex-row md:items-center md:justify-between gap-2 mb-6">
        <div>
          <h1 className="text-3xl font-bold text-transparent bg-clip-text bg-gradient-to-r from-sky-400 to-blue-500">
            Manage Beneficiaries
          </h1>
          <p className="text-sm text-gray-400 mt-1">
            Save trusted wallet IDs to send money quickly and safely.
          </p>
        </div>

        {wallet && (
          <div className="text-xs bg-white/5 border border-white/10 rounded-xl px-4 py-2 text-gray-300">
            <span className="block text-[10px] uppercase tracking-widest text-gray-500">
              Your Wallet
            </span>
            <span className="font-mono break-all">
              {wallet.wallet_id?.slice(0, 10)}...{wallet.wallet_id?.slice(-6)}
            </span>
          </div>
        )}
      </div>

      {/* Layout: Add form + list */}
      <div className="grid md:grid-cols-2 gap-6">
        {/* Add Beneficiary Card */}
        <div
          className="backdrop-blur-xl bg-white/5 border border-white/10
                     rounded-2xl p-6 shadow-lg"
        >
          <h2 className="text-lg font-semibold text-white mb-4">
            Add New Beneficiary
          </h2>

          <label className="text-xs text-gray-400 mb-1 block">
            Beneficiary Wallet ID
          </label>
          <input
            type="text"
            className="w-full p-3 rounded-xl bg-black/40 border border-gray-700
                       text-gray-200 text-sm focus:ring-2 focus:ring-sky-500 outline-none
                       placeholder:text-gray-500"
            placeholder="e.g. WALLET_xxx123..."
            value={newBeneficiary}
            onChange={(e) => setNewBeneficiary(e.target.value)}
          />

          <button
            onClick={addBeneficiary}
            disabled={loading}
            className="mt-4 w-full py-3 rounded-xl bg-sky-500 hover:bg-sky-600
                       text-white font-semibold text-sm tracking-wide
                       disabled:bg-sky-900/60 transition-colors"
          >
            {loading ? "Updating..." : "Save Beneficiary"}
          </button>

          <p className="mt-3 text-xs text-gray-500">
            Saved beneficiaries will appear on the{" "}
            <span className="text-sky-400">Send Money</span> screen for quick
            selection.
          </p>
        </div>

        {/* Beneficiaries List Card */}
        <div
          className="backdrop-blur-xl bg-white/5 border border-white/10
                     rounded-2xl p-6 shadow-lg"
        >
          <h2 className="text-lg font-semibold text-white mb-4">
            Your Beneficiaries
          </h2>

          {beneficiaries.length === 0 ? (
            <p className="text-gray-400 text-sm">
              You haven&apos;t added any beneficiaries yet.
              <br />
              Start by adding a wallet ID on the left.
            </p>
          ) : (
            <div className="space-y-3 max-h-72 overflow-y-auto pr-1 custom-scroll">
              {beneficiaries.map((b) => (
                <div
                  key={b}
                  className="flex items-center justify-between gap-3
                             bg-black/40 border border-gray-800 rounded-xl px-3 py-2"
                >
                  <div className="flex-1 min-w-0">
                    <p className="text-[11px] text-gray-500 uppercase tracking-wide">
                      Wallet ID
                    </p>
                    <p className="font-mono text-xs text-gray-200 break-all">
                      {b}
                    </p>
                  </div>

                  <button
                    onClick={() => removeBeneficiary(b)}
                    className="text-xs px-3 py-1 rounded-lg bg-red-500/80
                               hover:bg-red-600 text-white font-medium
                               shrink-0"
                  >
                    Remove
                  </button>
                </div>
              ))}
            </div>
          )}
        </div>
      </div>

      {/* Alerts */}
      {error && (
        <div className="mt-6 p-3 rounded-xl bg-red-900/40 border border-red-500/40 text-red-200 text-sm">
          {error}
        </div>
      )}

      {success && (
        <div className="mt-6 p-3 rounded-xl bg-emerald-900/40 border border-emerald-500/40 text-emerald-200 text-sm">
          {success}
        </div>
      )}
    </div>
  );
}
