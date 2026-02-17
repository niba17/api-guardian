import { useState } from "react";
import { Lock, User, ArrowRight, AlertCircle } from "lucide-react";
import { useNavigate } from "react-router-dom";
import { useAuth } from "../../../contexts/AuthContext";

export default function LoginForm() {
  const { login } = useAuth();
  const navigate = useNavigate();
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");
  const [form, setForm] = useState({ username: "", password: "" });

  const handleLogin = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError("");

    try {
      // Fetch ke Backend Go
      const response = await fetch("http://localhost:8080/api/login", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          username: form.username,
          password: form.password,
        }),
      });

      const data = await response.json();

      if (!response.ok) {
        throw new Error(data.error || "Login failed");
      }

      // Login Sukses
      login(data.user.username);
      localStorage.setItem("guardian_token", data.token); // Simpan Token
      navigate("/"); // Redirect ke Dashboard
    } catch (err) {
      console.error(err);
      // 👇 FIX: Mengganti (err: any) dengan pengecekan tipe yang aman
      if (err instanceof Error) {
        setError(err.message);
      } else {
        setError("Access Denied: Server Unreachable");
      }
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="w-full max-w-md bg-guardian-card border border-slate-800 p-8 rounded-3xl shadow-2xl relative overflow-hidden">
      <div className="absolute top-0 left-0 w-full h-1 bg-linear-to-r from-transparent via-guardian-primary to-transparent opacity-50"></div>

      <div className="text-center mb-8">
        <div className="inline-flex items-center justify-center w-16 h-16 rounded-full bg-guardian-primary/10 text-guardian-primary mb-4 animate-pulse">
          <Lock size={32} />
        </div>
        <h2 className="text-2xl font-bold text-white tracking-wide">
          SECURE ACCESS
        </h2>
        <p className="text-slate-500 text-sm mt-2">
          Identify yourself, Commander.
        </p>
      </div>

      <form onSubmit={handleLogin} className="space-y-6">
        {error && (
          <div className="bg-red-500/10 border border-red-500/20 text-red-400 p-3 rounded-xl flex items-center gap-2 text-sm">
            <AlertCircle size={16} /> {error}
          </div>
        )}

        <div className="space-y-4">
          <div className="relative group">
            <User
              className="absolute left-4 top-3.5 text-slate-500 group-focus-within:text-guardian-primary transition-colors"
              size={20}
            />
            <input
              type="text"
              placeholder="Username"
              className="w-full bg-slate-900/50 border border-slate-700 text-white pl-12 pr-4 py-3 rounded-xl focus:outline-none focus:border-guardian-primary focus:ring-1 focus:ring-guardian-primary transition-all placeholder:text-slate-600"
              value={form.username}
              onChange={(e) => setForm({ ...form, username: e.target.value })}
            />
          </div>

          <div className="relative group">
            <Lock
              className="absolute left-4 top-3.5 text-slate-500 group-focus-within:text-guardian-primary transition-colors"
              size={20}
            />
            <input
              type="password"
              placeholder="Password"
              className="w-full bg-slate-900/50 border border-slate-700 text-white pl-12 pr-4 py-3 rounded-xl focus:outline-none focus:border-guardian-primary focus:ring-1 focus:ring-guardian-primary transition-all placeholder:text-slate-600"
              value={form.password}
              onChange={(e) => setForm({ ...form, password: e.target.value })}
            />
          </div>
        </div>

        <button
          disabled={loading}
          className="w-full bg-guardian-primary hover:bg-emerald-600 text-white font-bold py-3 rounded-xl transition-all active:scale-95 flex items-center justify-center gap-2 disabled:opacity-50 disabled:cursor-not-allowed"
        >
          {loading ? (
            <span className="animate-spin h-5 w-5 border-2 border-white border-t-transparent rounded-full"></span>
          ) : (
            <>
              Authenticate <ArrowRight size={18} />
            </>
          )}
        </button>
      </form>

      <div className="mt-6 text-center">
        <p className="text-xs text-slate-600">
          Restricted Area. Unauthorized access is a federal offense.
        </p>
      </div>
    </div>
  );
}
