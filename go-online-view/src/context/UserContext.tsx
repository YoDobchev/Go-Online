import { createContext } from "react";

export type UserContextValue = {
    user: { email: string; username: string } | null;
    loading: boolean;
};

export const UserContext = createContext<UserContextValue>({
    user: null,
    loading: true,
});
