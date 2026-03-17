import { createContext } from "react";

export interface User {
  name: string;
  email: string;
  image: string;
  role: string;
}

export interface AuthContextType {
  user: User | null;
  loading: boolean;
}

export const AuthContext = createContext<AuthContextType>({
  user: null,
  loading: true,
});
