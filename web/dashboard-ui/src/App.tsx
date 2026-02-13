import {
  BrowserRouter as Router,
  Routes,
  Route,
  Navigate,
} from "react-router-dom";
import Dashboard from "./pages/Dashboard";
import Settings from "./pages/Settings";
import Threats from "./pages/Threats";
import Traffic from "./pages/Traffic";
import Login from "./pages/Login"; // 👈 Import Login
import ProtectedRoute from "./routes/ProtectedRoute"; // 👈 Import Satpam
import { AuthProvider } from "./contexts/AuthContext";

function App() {
  return (
    <AuthProvider>
      <Router>
        <Routes>
          <Route path="/login" element={<Login />} />
          <Route element={<ProtectedRoute />}>
            <Route path="/" element={<Dashboard />} />
            <Route path="/settings" element={<Settings />} />
            <Route path="/threats" element={<Threats />} />
            <Route path="/traffic" element={<Traffic />} />
          </Route>

          {/* 3. Catch-All (Kalau nyasar) */}
          <Route path="*" element={<Navigate to="/" replace />} />
        </Routes>
      </Router>
    </AuthProvider>
  );
}

export default App;
