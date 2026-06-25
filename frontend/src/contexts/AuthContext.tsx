import { createContext, useState, type ReactNode } from "react";
import type { Account } from "@/api";

interface AuthContextType {
    account: Account | null;
    setAccount: React.Dispatch<React.SetStateAction<Account | null>>;
}

export const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const AuthProvider = ({ children }: { children: ReactNode }) => {
    const [account, setAccount] = useState<Account | null>(null);

    return (
        <AuthContext.Provider
            value={{
                account,
                setAccount,
            }}
        >
            {children}
        </AuthContext.Provider>
    );
};
