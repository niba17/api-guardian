import axios from "axios";

const api = axios.create({
  baseURL: "http://localhost:8080/api", // Sesuaikan dengan port backend Bos
});

// 👇 TAMBAHKAN INI: Interceptor Request
api.interceptors.request.use(
  (config) => {
    // Ambil token dari localStorage
    const token = localStorage.getItem("guardian_token");

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

// Interceptor Response (Opsional: Handle kalau token expired)
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
