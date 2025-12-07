import { useEffect, useState } from "react";
import { apiGet } from "../api/api";

export default function SystemLogs() {
  const [logs, setLogs] = useState([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    async function load() {
      try {
        const res = await apiGet("/logs");
        if (res.success) setLogs(res.data || []);
      } catch (err) {
        console.error("Log fetch error:", err);
      }
      setLoading(false);
    }

    load();
  }, []);

  if (loading)
    return (
      <p className="text-gray-400 p-6 text-center animate-pulse">
        Loading logs...
      </p>
    );

  if (logs.length === 0)
    return (
      <p className="p-6 text-gray-400 text-center">
        No system logs available.
      </p>
    );

  // Badge color by level
  const badgeColor = (level) => {
    switch (level) {
      case "success":
        return "bg-green-700/40 border border-green-500 text-green-300";
      case "error":
        return "bg-red-700/40 border border-red-500 text-red-300";
      default:
        return "bg-blue-700/40 border border-blue-500 text-blue-300";
    }
  };

  return (
    <div className="max-w-5xl mx-auto p-8 animate-fadeIn">

      {/* Header */}
      <h1 className="text-4xl font-bold mb-8 text-transparent bg-clip-text 
        bg-gradient-to-r from-purple-400 via-pink-400 to-blue-400 drop-shadow-lg">
        System Logs
      </h1>

      {/* Log Container */}
      <div className="space-y-5">
        {logs.map((log, i) => (
          <div
            key={i}
            className="
              backdrop-blur-xl bg-white/5 border border-white/10
              p-5 rounded-2xl shadow-lg hover:shadow-purple-500/20 
              transition duration-200
              text-gray-200
            "
          >
            {/* Header line: timestamp + badge */}
            <div className="flex justify-between items-center">
              <p className="text-sm text-gray-400">
                {new Date(log.timestamp * 1000).toLocaleString()}
              </p>

              <span
                className={`${badgeColor(
                  log.level
                )} px-3 py-1 rounded-full text-xs font-semibold`}
              >
                {log.level.toUpperCase()}
              </span>
            </div>

            {/* Message */}
            <p className="mt-3 text-lg leading-relaxed text-gray-100">
              {log.message}
            </p>
          </div>
        ))}
      </div>
    </div>
  );
}
