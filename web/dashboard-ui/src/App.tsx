import {
  BrowserRouter as Router,
  Routes,
  Route,
  Navigate,
} from "react-router-dom";
import Overview from "./pages/Overview";
import Config from "./pages/Config";
import WAF from "./pages/WAF";
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
            <Route path="/" element={<Overview />} />
            <Route path="/settings" element={<Config />} />
            <Route path="/threats" element={<WAF />} />
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
