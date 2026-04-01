import { useEffect, useState, useRef } from "react";
import {
  Bell,
  ChefHat,
  CheckCircle2,
  ArrowLeft,
  Loader2,
  Volume2,
  VolumeX,
} from "lucide-react";
import axios from "axios";
import toast from "react-hot-toast";
import { Link } from "react-router-dom";
import { restaurantServiceUrl } from "@/lib/config";

interface LiveOrder {
  order_id: string;
  grand_total: number;
  message: string;
  items: { name: string; quantity: number }[];
}

export default function RestaurantDashboard() {
  const [liveOrders, setLiveOrders] = useState<LiveOrder[]>([]);
  const [isConnected, setIsConnected] = useState(false);
  const [processingId, setProcessingId] = useState<string | null>(null);

  // ✅ Proper State Management
  const [restaurantId, setRestaurantId] = useState<string | null>(null);
  const [isInitializing, setIsInitializing] = useState(true);
  const [audioUnlocked, setAudioUnlocked] = useState(false);

  const wsRef = useRef<WebSocket | null>(null);
  const audioRef = useRef<HTMLAudioElement | null>(null);

  // 1. Initialize Audio Object
  useEffect(() => {
    audioRef.current = new Audio(
      "https://assets.mixkit.co/active_storage/sfx/2869/2869-preview.mp3",
    );
  }, []);

  // 2. Fetch the correct RESTAURANT ID (Not the User ID!)
  useEffect(() => {
    const fetchRestaurantProfile = async () => {
      try {
        const res = await axios.get(`${restaurantServiceUrl}/restaurant/read`, {
          withCredentials: true,
        });
        const data = res.data?.data ?? res.data;
        if (data?.id) {
          setRestaurantId(data.id);
        } else {
          toast.error("Restaurant profile not found.");
        }
      } catch (error) {
        console.error("Failed to fetch restaurant:", error);
      } finally {
        setIsInitializing(false);
      }
    };
    fetchRestaurantProfile();
  }, []);

  // 3. Connect WebSocket & Fetch Orders ONLY after we have the Restaurant ID
  useEffect(() => {
    if (!restaurantId) return;

    let ws: WebSocket;
    let isComponentMounted = true;

    const fetchInitialOrders = async () => {
      try {
        const res = await axios.get(
          `${restaurantServiceUrl}/restaurant/${restaurantId}/orders/active`,
          { withCredentials: true },
        );

        const fetchedOrders = (
          res.data?.data?.orders ||
          res.data?.orders ||
          []
        ).map((o: any) => ({
          order_id: o.id,
          grand_total: o.grand_total,
          message: `Status: ${o.status.toUpperCase()}`,
          items: o.items,
        }));

        if (isComponentMounted) setLiveOrders(fetchedOrders);
      } catch (error) {
        console.error("Failed to fetch initial orders:", error);
      }
    };

    const connectWebSocket = () => {
      ws = new WebSocket(`ws://localhost:8003/ws/restaurant/${restaurantId}`);

      ws.onopen = () => {
        if (!isComponentMounted) return;
        console.log("🟢 Connected to Realtime Dashboard");
        setIsConnected(true);
      };

      ws.onmessage = (event) => {
        if (!isComponentMounted) return;
        const newOrder = JSON.parse(event.data);
        console.log("🚨 NEW ORDER ARRIVED:", newOrder);

        setLiveOrders((prev) => {
          if (prev.some((o) => o.order_id === newOrder.order_id)) return prev;
          return [newOrder, ...prev];
        });

        // ✅ Play the sound safely
        if (audioRef.current && audioUnlocked) {
          audioRef.current.currentTime = 0; // Reset to start
          audioRef.current
            .play()
            .catch((e) => console.log("Audio still blocked:", e));
        }
      };

      ws.onclose = () => {
        if (!isComponentMounted) return;
        console.log("🔴 Disconnected from Realtime Dashboard");
        setIsConnected(false);
        setTimeout(connectWebSocket, 3000);
      };
    };

    fetchInitialOrders().then(() => {
      if (isComponentMounted) connectWebSocket();
    });

    return () => {
      isComponentMounted = false;
      if (ws) {
        ws.onclose = null;
        ws.close();
      }
    };
  }, [restaurantId, audioUnlocked]); // Re-run if audio unlocks so the WS closure gets the new state

  // ✅ 4. The Audio Unlocker Function
  const handleUnlockAudio = () => {
    if (audioRef.current) {
      // We play it silently first to satisfy the browser's interaction policy
      audioRef.current.volume = 0;
      audioRef.current
        .play()
        .then(() => {
          audioRef.current!.pause();
          audioRef.current!.currentTime = 0;
          audioRef.current!.volume = 1; // Turn volume back up for real orders!
          setAudioUnlocked(true);
          toast.success("Audio alerts enabled!");
        })
        .catch((e) => console.error("Could not unlock audio", e));
    }
  };

  const handleMarkPreparing = async (orderId: string) => {
    try {
      setProcessingId(orderId);
      await axios.patch(
        `${restaurantServiceUrl}/order/${orderId}/status`,
        {
          status: "preparing",
          restaurant_id: restaurantId,
        },
        { withCredentials: true },
      );
      toast.success("Order moved to kitchen!");
      setLiveOrders((prev) => prev.filter((o) => o.order_id !== orderId));
    } catch (error) {
      toast.error("Failed to update order status");
    } finally {
      setProcessingId(null);
    }
  };

  if (isInitializing) {
    return (
      <div className="min-h-screen bg-gray-900 flex flex-col items-center justify-center text-white">
        <Loader2 className="w-10 h-10 text-red-500 animate-spin mb-4" />
        <p>Loading your kitchen display...</p>
      </div>
    );
  }

  if (!restaurantId && !isInitializing) {
    return (
      <div className="min-h-screen bg-gray-900 flex flex-col items-center justify-center text-white">
        <ChefHat className="w-16 h-16 text-red-500 mb-4" />
        <h2 className="text-2xl font-bold mb-2">No Restaurant Found</h2>
        <p className="text-gray-400 mb-6">
          You need to set up a restaurant profile first.
        </p>
        <Link
          to="/"
          className="bg-red-500 hover:bg-red-600 px-6 py-2 rounded-lg font-bold"
        >
          Go Back
        </Link>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-900 text-gray-100 p-8">
      <div className="max-w-6xl mx-auto">
        {/* Header */}
        <div className="flex items-center justify-between mb-8 pb-4 border-b border-gray-800">
          <div className="flex items-center gap-4">
            <Link
              to="/"
              className="p-2 hover:bg-gray-800 rounded-lg transition"
            >
              <ArrowLeft className="w-6 h-6 text-gray-400" />
            </Link>
            <ChefHat className="w-8 h-8 text-red-500" />
            <h1 className="text-2xl font-bold text-white">
              Kitchen Display System
            </h1>
          </div>

          <div className="flex items-center gap-4">
            {/* 🔊 Audio Unlock Button */}
            {!audioUnlocked ? (
              <button
                onClick={handleUnlockAudio}
                className="flex items-center gap-2 bg-yellow-500/20 text-yellow-500 border border-yellow-500/50 hover:bg-yellow-500/30 px-4 py-1.5 rounded-full text-sm font-bold transition animate-pulse"
              >
                <VolumeX className="w-4 h-4" /> Click to Enable Sounds
              </button>
            ) : (
              <div className="flex items-center gap-2 text-gray-500 px-4 py-1.5 rounded-full text-sm border border-gray-800">
                <Volume2 className="w-4 h-4" /> Sounds Active
              </div>
            )}

            <div
              className={`flex items-center gap-2 px-3 py-1.5 rounded-full text-sm font-medium ${isConnected ? "bg-green-500/10 text-green-400" : "bg-red-500/10 text-red-400"}`}
            >
              <div
                className={`w-2 h-2 rounded-full ${isConnected ? "bg-green-400 animate-pulse" : "bg-red-500"}`}
              />
              {isConnected ? "Live & Listening" : "Disconnected"}
            </div>
          </div>
        </div>

        {/* Orders Grid */}
        {liveOrders.length === 0 ? (
          <div className="flex flex-col items-center justify-center h-[50vh] text-gray-500">
            <Bell className="w-16 h-16 mb-4 opacity-20" />
            <h2 className="text-xl font-semibold">Waiting for new orders...</h2>
            <p className="text-sm mt-2">
              They will appear here instantly when paid.
            </p>
          </div>
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {liveOrders.map((order) => (
              <div
                key={order.order_id}
                className="bg-gray-800 rounded-2xl p-6 border border-gray-700 shadow-xl animate-in slide-in-from-bottom-4 duration-500"
              >
                <div className="flex justify-between items-start mb-4">
                  <div>
                    <span className="inline-block bg-red-500 text-white text-xs font-bold px-2 py-1 rounded-md mb-2">
                      NEW ORDER
                    </span>
                    <p className="text-gray-400 text-xs font-mono">
                      ID: {order.order_id.slice(-6).toUpperCase()}
                    </p>
                  </div>
                  <p className="font-bold text-lg text-green-400">
                    ₹{order.grand_total}
                  </p>
                </div>

                <div className="space-y-3 mb-6">
                  {order.items.map((item, idx) => (
                    <div key={idx} className="flex justify-between text-sm">
                      <span className="font-medium text-gray-200">
                        <span className="text-gray-500 mr-2">
                          {item.quantity}x
                        </span>{" "}
                        {item.name}
                      </span>
                    </div>
                  ))}
                </div>

                <button
                  onClick={() => handleMarkPreparing(order.order_id)}
                  disabled={processingId === order.order_id}
                  className="w-full flex items-center justify-center gap-2 bg-gray-700 hover:bg-green-600 text-white py-3 rounded-xl font-bold transition-colors disabled:opacity-50"
                >
                  <CheckCircle2 className="w-5 h-5" />
                  {processingId === order.order_id
                    ? "Updating..."
                    : "Mark Preparing"}
                </button>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  );
}
