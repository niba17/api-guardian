import LoginForm from "../features/auth/components/LoginForm";

export default function Login() {
  return (
    <div className="min-h-screen bg-guardian-dark flex items-center justify-center p-4 relative overflow-hidden">
      {/* Background Grid Pattern (Opsional biar keren) */}
      <div className="absolute inset-0 bg-[linear-gradient(to_right,#80808012_1px,transparent_1px),linear-gradient(to_bottom,#80808012_1px,transparent_1px)] bg-size-[24px_24px]"></div>

      {/* Gradient Glow */}
      <div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-125 h-125 bg-guardian-primary/20 rounded-full blur-[100px] opacity-20 pointer-events-none"></div>

      <div className="relative z-10 w-full flex justify-center">
        <LoginForm />
      </div>
    </div>
  );
}
