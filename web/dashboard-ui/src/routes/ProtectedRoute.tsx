import { Navigate, Outlet } from "react-router-dom";
import { useAuth } from "../contexts/AuthContext"; // 👈 Pakai hook

export default function ProtectedRoute() {
  const { isAuthenticated } = useAuth();

  if (!isAuthenticated) {
    return <Navigate to="/login" replace />;
  }

  return <Outlet />;
}
