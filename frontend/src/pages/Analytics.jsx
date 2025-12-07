import { useEffect, useState } from "react";
import { apiGet } from "../api/api";
import {
  ArrowPathIcon,
  CheckCircleIcon,
  BanknotesIcon,
} from "@heroicons/react/24/solid";

export default function Analytics() {
  const [data, setData] = useState(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    async function load() {
      const res = await apiGet("/analytics/system");
      if (res.success) setData(res.data);
      setLoading(false);
    }
    load();
  }, []);

  if (loading)
    return (
      <p className="text-center text-gray-400 mt-10 animate-pulse">
        Loading analytics...
      </p>
    );

  if (!data)
    return <p className="p-6 text-center text-red-500">Failed to load.</p>;

  return (
    <div className="max-w-3xl mx-auto p-8 animate-fadeIn">

      {/* Title */}
      <h1
        className="text-4xl font-bold mb-10 text-transparent bg-clip-text 
        bg-gradient-to-r from-green-300 via-blue-300 to-cyan-400 text-center"
      >
        System Analytics
      </h1>

      {/* Vertical Stats */}
      <div className="space-y-6">

        <AnalyticsCard
          icon={<ArrowPathIcon className="w-10 h-10 text-yellow-400" />}
          label="Pending Transactions"
          value={data.pending_transactions}
          color="from-yellow-500 to-amber-600"
        />

        <AnalyticsCard
          icon={<CheckCircleIcon className="w-10 h-10 text-blue-400" />}
          label="Completed Transactions"
          value={data.completed_transactions}
          color="from-blue-500 to-purple-600"
        />

        <AnalyticsCard
          icon={<BanknotesIcon className="w-10 h-10 text-green-400" />}
          label="Total Zakat Collected"
          value={`${data.total_zakat_collected} ₿`}
          color="from-green-500 to-emerald-600"
        />

      </div>
    </div>
  );
}

function AnalyticsCard({ icon, label, value, color }) {
  return (
    <div
      className="
        flex items-center gap-6 p-6 rounded-2xl
        backdrop-blur-xl bg-white/5 border border-white/10
        shadow-lg hover:shadow-2xl hover:shadow-green-500/20
        transition duration-300 transform hover:-translate-y-1
      "
    >
      {/* Icon Section */}
      <div className="p-4 rounded-xl bg-black/30 border border-white/5 shadow-inner">
        {icon}
      </div>

      {/* Text Section */}
      <div>
        <p className="text-gray-300 text-sm">{label}</p>

        <p
          className={`text-4xl font-bold mt-1 bg-gradient-to-r ${color}
          bg-clip-text text-transparent drop-shadow-lg`}
        >
          {value}
        </p>
      </div>
    </div>
  );
}
