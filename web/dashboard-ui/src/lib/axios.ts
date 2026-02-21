import axios from "axios";

const api = axios.create({
  baseURL: "http://localhost:8080/api", // Sesuaikan dengan port backend Bos
});

// Interceptor Request
api.interceptors.request.use(
  (config) => {
    // Ambil token dari localStorage
    const token = localStorage.getItem("guardian_token");

    // 🚀 FIX: Suntikkan Kunci Gerbang Depan (API Key) ke semua request!
    config.headers["X-API-KEY"] = "kunci-rahasia-bos-123";

    if (token) {
      // Pasang di header Authorization: Bearer <token>
      config.headers.Authorization = `Bearer ${token}`;
    }

    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// Interceptor Response
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      // Jika dapat 401 dari server, paksa logout atau balik ke login
      console.error("Session Expired or Invalid Token");
      localStorage.removeItem("guardian_token");
      window.location.href = "/login";
    }
    return Promise.reject(error);
  }
);

export default api;
