import React from "react";
import ReactDOM from "react-dom/client";
import App from "./App.jsx";
import "./index.css";

import {
  BrowserRouter as Router,
  Routes,
  Route,
} from "react-router-dom";

// Pages
import Dashboard from "./pages/Dashboard.jsx";
import SendMoney from "./pages/SendMoney.jsx";
import BlockExplorer from "./pages/BlockExplorer.jsx";
import SystemLogs from "./pages/SystemLogs.jsx";
import TransactionHistory from "./pages/TransactionHistory.jsx";
import Register from "./pages/Register.jsx";
import Login from "./pages/Login.jsx";
import Beneficiaries from "./pages/Beneficiaries.jsx";
import Analytics from "./pages/Analytics";
import ZakatReport from "./pages/ZakatReport";
import Profile from "./pages/Profile.jsx";
import Reports from "./pages/Reports.jsx";

// Protected Route
import ProtectedRoute from "./components/ProtectedRoute.jsx";

ReactDOM.createRoot(document.getElementById("root")).render(
  <Router>
    <Routes>

      {/* ------------------------- */}
      {/*  PUBLIC FULLSCREEN ROUTES */}
      {/* ------------------------- */}
      <Route path="/login" element={<Login />} />
      <Route path="/register" element={<Register />} />

      {/* ------------------------- */}
      {/*      APP LAYOUT ROUTES    */}
      {/* ------------------------- */}
      <Route path="/" element={<App />}>

        <Route
          index
          element={
            <ProtectedRoute>
              <Dashboard />
            </ProtectedRoute>
          }
        />

        <Route
          path="profile"
          element={
            <ProtectedRoute>
              <Profile />
            </ProtectedRoute>
          }
        />

        <Route
          path="send"
          element={
            <ProtectedRoute>
              <SendMoney />
            </ProtectedRoute>
          }
        />

        <Route
          path="explorer"
          element={
            <ProtectedRoute>
              <BlockExplorer />
            </ProtectedRoute>
          }
        />

        <Route
          path="logs"
          element={
            <ProtectedRoute>
              <SystemLogs />
            </ProtectedRoute>
          }
        />

        <Route
          path="history/:walletId"
          element={
            <ProtectedRoute>
              <TransactionHistory />
            </ProtectedRoute>
          }
        />

        <Route
          path="beneficiaries"
          element={
            <ProtectedRoute>
              <Beneficiaries />
            </ProtectedRoute>
          }
        />

        <Route
          path="analytics"
          element={
            <ProtectedRoute>
              <Analytics />
            </ProtectedRoute>
          }
        />

        <Route
          path="reports"
          element={
            <ProtectedRoute>
              <Reports />
            </ProtectedRoute>
          }
        />

        <Route
          path="zakat-report"
          element={
            <ProtectedRoute>
              <ZakatReport />
            </ProtectedRoute>
          }
        />

      </Route>

    </Routes>
  </Router>
);
