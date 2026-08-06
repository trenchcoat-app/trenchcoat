import { createContext } from "react";
import type { Account } from "@/api";

export interface AuthContextType {
    account: Account | null;
    setAccount: React.Dispatch<React.SetStateAction<Account | null>>;
    isAuthenticated: boolean;
}

export const AuthContext = createContext<AuthContextType | undefined>(
    undefined,
);