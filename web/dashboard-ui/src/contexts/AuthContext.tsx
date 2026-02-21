import { createContext, useContext, useState, type ReactNode } from "react";

interface User {
  username: string;
  role: string;
}

interface AuthContextType {
  user: User | null;
  // 🚀 FIX 1: Fungsi login sekarang menerima token asli dari luar
  login: (username: string, token: string) => void;
  logout: () => void;
  isAuthenticated: boolean;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<User | null>(() => {
    const token = localStorage.getItem("guardian_token");
    if (token) return { username: "admin", role: "Commander" };
    return null;
  });

  // 🚀 FIX 2: Simpan token yang dikirimkan, bukan string "rahasia-negara"
  const login = (username: string, token: string) => {
    localStorage.setItem("guardian_token", token);
    setUser({ username, role: "Commander" });
  };

  const logout = () => {
    localStorage.removeItem("guardian_token");
    setUser(null);
  };

  return (
    <AuthContext.Provider
      value={{ user, login, logout, isAuthenticated: !!user }}
    >
      {children}
    </AuthContext.Provider>
  );
}

// eslint-disable-next-line react-refresh/only-export-components
export function useAuth() {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error("useAuth must be used within an AuthProvider");
  }
  return context;
}
