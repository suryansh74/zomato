export default function Login() {
  const handleLogin = () => {
    window.location.href = "http://localhost:8000/api/auth/login";
  };

  return <button onClick={handleLogin}>Login with Google</button>;
}
