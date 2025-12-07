import { useEffect, useState } from "react";
import { apiGet, apiPost } from "../api/api";

export default function Profile() {
  const walletId = localStorage.getItem("walletId");

  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [message, setMessage] = useState("");

  // Profile fields
  const [name, setName] = useState("");
  const [email, setEmail] = useState("");
  const [cnic, setCnic] = useState("");

  // Wallet details
  const [publicKey, setPublicKey] = useState("");
  const [beneficiaries, setBeneficiaries] = useState([]);

  // =====================================================
  // LOAD PROFILE DATA
  // =====================================================
  useEffect(() => {
    async function loadProfile() {
      try {
        const res = await apiGet(`/user/profile/${walletId}`);
        if (!res.success) {
          setMessage("Failed to load profile");
          setLoading(false);
          return;
        }

        const data = res.data;

        setName(data.name || "");
        setEmail(data.email || "");
        setCnic(data.cnic || "");
        setPublicKey(data.public_key || "");
        setBeneficiaries(data.beneficiaries || []);

      } catch {
        setMessage("Error loading profile");
      }
      setLoading(false);
    }

    loadProfile();
  }, [walletId]);

  // =====================================================
  // UPDATE PROFILE
  // =====================================================
  async function handleSave() {
    setSaving(true);
    setMessage("");

    try {
      const res = await apiPost("/user/profile/update", {
        wallet_id: walletId,
        name,
        email,
        cnic,
      });

      if (!res.success) {
        setMessage(res.message || "Update failed");
      } else {
        setMessage("✔ Profile updated successfully!");
      }

    } catch {
      setMessage("Failed to update profile");
    }

    setSaving(false);
  }

  if (loading)
    return (
      <p className="text-center text-gray-300 mt-16 animate-pulse">
        Loading profile...
      </p>
    );

  // =====================================================
  // BEAUTIFUL WEB3 UI
  // =====================================================
  return (
    <div className="max-w-5xl mx-auto py-10 space-y-10">

      <h1 className="text-4xl font-bold text-blue-400 drop-shadow mb-6">
        Profile Settings
      </h1>

      {/* PERSONAL INFO SECTION */}
      <div className="
        bg-white/5 backdrop-blur-xl border border-white/10
        p-8 rounded-2xl shadow-lg animate-fadeIn
      ">
        <h2 className="text-2xl font-semibold text-white mb-6">
          Personal Information
        </h2>

        <div className="space-y-5">

          {/* NAME */}
          <div>
            <label className="text-gray-300 text-sm">Full Name</label>
            <input
              type="text"
              className="
                w-full mt-1 p-3 rounded-lg bg-black/40 border border-gray-700 
                text-gray-200 focus:ring-2 focus:ring-blue-500 outline-none
              "
              value={name}
              onChange={(e) => setName(e.target.value)}
            />
          </div>

          {/* EMAIL */}
          <div>
            <label className="text-gray-300 text-sm">Email</label>
            <input
              type="email"
              className="
                w-full mt-1 p-3 rounded-lg bg-black/40 border border-gray-700 
                text-gray-200 focus:ring-2 focus:ring-blue-500 outline-none
              "
              value={email}
              onChange={(e) => setEmail(e.target.value)}
            />
          </div>

          {/* CNIC */}
          <div>
            <label className="text-gray-300 text-sm">CNIC</label>
            <input
              type="text"
              className="
                w-full mt-1 p-3 rounded-lg bg-black/40 border border-gray-700 
                text-gray-200 focus:ring-2 focus:ring-blue-500 outline-none
              "
              value={cnic}
              onChange={(e) => setCnic(e.target.value)}
            />
          </div>

          {/* SAVE BUTTON */}
          <button
            onClick={handleSave}
            disabled={saving}
            className="
              mt-4 px-5 py-3 w-full rounded-lg 
              bg-blue-600 hover:bg-blue-700 
              text-white font-semibold 
              shadow-md hover:shadow-blue-500/30
              transition
            "
          >
            {saving ? "Saving..." : "Save Changes"}
          </button>

          {/* SUCCESS / ERROR MESSAGE */}
          {message && (
            <div className="
              mt-4 p-3 rounded-lg 
              bg-green-900/40 text-green-300 text-center
            ">
              {message}
            </div>
          )}
        </div>
      </div>

      {/* WALLET SECTION */}
      <div className="
        bg-white/5 backdrop-blur-xl border border-white/10
        p-8 rounded-2xl shadow-lg animate-fadeIn
      ">
        <h2 className="text-2xl font-semibold text-white mb-6">
          Wallet Details
        </h2>

        {/* Wallet ID */}
        <div className="mb-6">
          <label className="text-gray-300 text-sm">Wallet ID</label>
          <input
            className="
              w-full mt-1 p-3 rounded-lg bg-black/40 border border-gray-700 
              text-gray-300 outline-none
            "
            value={walletId}
            readOnly
          />
        </div>

        {/* Public Key */}
        <div className="mb-6">
          <label className="text-gray-300 text-sm">Public Key</label>
          <textarea
            className="
              w-full mt-1 p-3 rounded-lg bg-black/40 border border-gray-700 
              text-gray-300 outline-none
            "
            rows={4}
            value={publicKey}
            readOnly
          />
        </div>

        {/* Beneficiaries */}
        <div>
          <label className="text-gray-300 text-sm">Beneficiaries</label>

          <div className="mt-2 flex flex-wrap gap-2">
            {beneficiaries.length === 0 ? (
              <p className="text-gray-400 text-sm">
                No beneficiaries added
              </p>
            ) : (
              beneficiaries.map((b, i) => (
                <span
                  key={i}
                  className="
                    px-3 py-1 rounded-full 
                    bg-blue-900/40 border border-blue-700
                    text-blue-300 text-sm break-all
                  "
                >
                  {b}
                </span>
              ))
            )}
          </div>
        </div>
      </div>
    </div>
  );
}
