import { createContext, useEffect, useState, type ReactNode } from "react";
import type { Account } from "@/api";

export interface AuthContextType {
    account: Account | null;
    setAccount: React.Dispatch<React.SetStateAction<Account | null>>;
    isAuthenticated: boolean;
}

export const AuthContext = createContext<AuthContextType | undefined>(undefined);

const STORAGE_KEY = "trenchcoat_account";

export const AuthProvider = ({ children }: { children: ReactNode }) => {
    const [account, setAccount] = useState<Account | null>(() => {
        const storedAccount = localStorage.getItem(STORAGE_KEY);
        try {
            return storedAccount ? (JSON.parse(storedAccount) as Account) : null;
        } catch {
            return null;
        }
    });

    useEffect(() => {
        if (account) {
            localStorage.setItem(STORAGE_KEY, JSON.stringify(account));
        } else {
            localStorage.removeItem(STORAGE_KEY);
        }
    }, [account]);

    return (
        <AuthContext.Provider
            value={{
                account,
                setAccount,
                isAuthenticated: !!account,
            }}
        >
            {children}
        </AuthContext.Provider>
    );
};
