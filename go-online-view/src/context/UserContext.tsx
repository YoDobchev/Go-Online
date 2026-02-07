import { createContext } from "react";

export type UserContextValue = {
    user: { email: string; username: string; isInGameWithID: string } | null;
    loading: boolean;
    setUser: (
        user: {
            email: string;
            username: string;
            isInGameWithID: string;
        } | null,
    ) => void;
};

export const UserContext = createContext<UserContextValue>({
    user: null,
    loading: true,
    setUser: () => {},
});
