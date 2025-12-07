import { Link, useNavigate, useLocation } from "react-router-dom";

export default function Navbar() {
  const navigate = useNavigate();
  const location = useLocation();

  const token = localStorage.getItem("token");
  const walletId = localStorage.getItem("walletId");

  function handleLogout() {
    localStorage.clear();
    navigate("/login");
  }

  // highlight active route
  const isActive = (path) =>
    location.pathname === path ||
    (path.includes("history") && location.pathname.includes("history"));

  return (
    <nav className="
      w-full px-8 py-4
      bg-black/40 backdrop-blur-xl 
      border-b border-white/10
      sticky top-0 z-50 
      shadow-lg 
      flex justify-between items-center
    ">
      {/* BRAND */}
      <Link
        to="/"
        className="text-2xl font-bold text-blue-400 drop-shadow-lg hover:text-blue-300 transition"
      >
        Blockchain Wallet
      </Link>

      {/* NAV LINKS */}
      <div className="flex items-center space-x-6 text-gray-300 font-medium">

        {!token ? (
          <>
            <Link
              to="/login"
              className="hover:text-white transition"
            >
              Login
            </Link>
            <Link
              to="/register"
              className="hover:text-white transition"
            >
              Register
            </Link>
          </>
        ) : (
          <>
            {[
              { label: "Dashboard", path: "/" },
              { label: "Profile", path: "/profile" },
              { label: "Send", path: "/send" },
              { label: "Explorer", path: "/explorer" },
              { label: "Logs", path: "/logs" },
              { label: "Analytics", path: "/analytics" },
              { label: "Reports", path: "/reports" },
              { label: "Beneficiaries", path: "/beneficiaries" },
              { label: "Zakat Report", path: "/zakat-report" }
            ].map((item, i) => (
              <Link
                key={i}
                to={item.path}
                className={`
                  px-2 py-1 rounded-md transition
                  ${isActive(item.path)
                    ? "text-white bg-white/10 shadow-inner"
                    : "hover:text-white hover:bg-white/5"}
                `}
              >
                {item.label}
              </Link>
            ))}

            {walletId && (
              <Link
                to={`/history/${walletId}`}
                className={`
                  px-2 py-1 rounded-md transition
                  ${isActive(`/history/${walletId}`)
                    ? "text-white bg-white/10 shadow-inner"
                    : "hover:text-white hover:bg-white/5"}
                `}
              >
                History
              </Link>
            )}

            {/* LOGOUT BUTTON */}
            <button
              onClick={handleLogout}
              className="
                ml-4 px-4 py-2 rounded-lg 
                bg-red-600/80 hover:bg-red-600 
                text-white font-semibold 
                shadow-md hover:shadow-red-500/30
                transition
              "
            >
              Logout
            </button>
          </>
        )}
      </div>
    </nav>
  );
}
