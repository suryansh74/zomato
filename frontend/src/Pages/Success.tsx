import { Link, useSearchParams } from "react-router-dom";
import { CheckCircle } from "lucide-react";
import { useEffect } from "react";

export default function Success() {
  const [searchParams] = useSearchParams();
  const sessionId = searchParams.get("session_id");

  useEffect(() => {
    // If you want to force the cart to clear visually for now,
    // you can just reload the page once when they click "Back to Home"
    console.log("Stripe Session ID:", sessionId);
  }, [sessionId]);

  return (
    <div className="min-h-[80vh] flex flex-col items-center justify-center bg-gray-50 px-4">
      <div className="bg-white p-8 rounded-3xl shadow-sm border border-gray-100 text-center max-w-md w-full">
        <div className="flex justify-center mb-6">
          <CheckCircle className="w-24 h-24 text-green-500" />
        </div>

        <h1 className="text-3xl font-bold text-gray-900 mb-2">
          Payment Successful!
        </h1>
        <p className="text-gray-500 mb-8">
          Your order has been placed and sent to the restaurant. They will start
          preparing it shortly.
        </p>

        <div className="space-y-3">
          <Link
            to="/"
            onClick={() => (window.location.href = "/")} // Force a hard refresh to clear the cart UI state
            className="block w-full bg-red-500 text-white font-bold py-3 rounded-xl hover:bg-red-600 transition"
          >
            Back to Home
          </Link>
        </div>
      </div>
    </div>
  );
}
