import { useEffect, useState } from "react";
import { apiGet } from "../api/api";

export default function Reports() {
  const walletId = localStorage.getItem("walletId");

  const [loading, setLoading] = useState(true);
  const [report, setReport] = useState(null);
  const [error, setError] = useState("");

  useEffect(() => {
    async function loadReport() {
      try {
        const res = await apiGet(`/reports/wallet/${walletId}`);

        if (!res.success) {
          setError("Failed to load reports");
          setLoading(false);
          return;
        }

        setReport(res.data);
      } catch (err) {
        console.error(err);
        setError("Backend error while loading reports");
      }

      setLoading(false);
    }

    loadReport();
  }, [walletId]);

  if (loading)
    return (
      <p className="text-center text-gray-400 mt-10 animate-pulse">
        Loading Reports...
      </p>
    );

  if (error)
    return <p className="text-center text-red-500 mt-10">{error}</p>;

  const {
    totalSent,
    totalReceived,
    zakatPaid,
    monthlySummary,
    transactionCount,
    receivedCount,
    sentCount,
    blocksMined,
  } = report;

  return (
    <div className="max-w-6xl mx-auto p-8 animate-fadeIn">

      {/* Title */}
      <h1
        className="text-4xl font-bold mb-10 text-transparent bg-clip-text 
        bg-gradient-to-r from-purple-400 via-pink-300 to-red-400"
      >
        Wallet Reports
      </h1>

      {/* Stats Grid */}
      <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-6">

        <ReportCard
          label="Total Sent"
          value={`${totalSent} coins`}
          sub={`${sentCount} transactions`}
          color="from-red-500 to-rose-600"
        />

        <ReportCard
          label="Total Received"
          value={`${totalReceived} coins`}
          sub={`${receivedCount} transactions`}
          color="from-green-500 to-emerald-600"
        />

        <ReportCard
          label="Total Zakat Deducted"
          value={`${zakatPaid} coins`}
          color="from-blue-500 to-indigo-600"
        />

        <ReportCard
          label="Blocks Mined"
          value={blocksMined}
          color="from-purple-500 to-fuchsia-600"
        />

        <ReportCard
          label="Total Transactions"
          value={transactionCount}
          color="from-orange-500 to-amber-600"
        />
      </div>

      {/* Monthly Summary Table */}
      <div className="mt-10">
        <h2 className="text-3xl font-semibold mb-4 text-white">
          Monthly Summary
        </h2>

        <div
          className="overflow-x-auto backdrop-blur-xl bg-white/5 border border-white/10 
          rounded-2xl shadow-xl"
        >
          <table className="min-w-full text-gray-200">
            <thead className="bg-white/10 text-gray-300">
              <tr>
                <th className="p-3 text-left">Month</th>
                <th className="p-3 text-left">Sent</th>
                <th className="p-3 text-left">Received</th>
                <th className="p-3 text-left">Zakat</th>
                <th className="p-3 text-left">Transactions</th>
              </tr>
            </thead>

            <tbody>
              {monthlySummary.map((row, i) => (
                <tr
                  key={i}
                  className="border-t border-white/10 hover:bg-white/5 transition"
                >
                  <td className="p-3">{row.month}</td>
                  <td className="p-3 text-red-400">{row.sent}</td>
                  <td className="p-3 text-green-400">{row.received}</td>
                  <td className="p-3 text-blue-400">{row.zakat}</td>
                  <td className="p-3">{row.txCount}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
}

function ReportCard({ label, value, sub, color }) {
  return (
    <div
      className="
        backdrop-blur-xl bg-white/5 border border-white/10
        p-6 rounded-2xl shadow-lg hover:shadow-2xl hover:shadow-purple-500/20
        transition duration-200 transform hover:-translate-y-1
      "
    >
      <p className="text-sm text-gray-300 mb-1">{label}</p>

      <p
        className={`
          text-3xl font-bold bg-gradient-to-r ${color}
          text-transparent bg-clip-text mb-1
        `}
      >
        {value}
      </p>

      {sub && <p className="text-gray-400 text-sm">{sub}</p>}
    </div>
  );
}
