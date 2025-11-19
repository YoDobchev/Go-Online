import { UserProvider } from "./context/UserProvider";
import { BrowserRouter } from "react-router-dom";

export function AppProviders({ children }: { children: React.ReactNode }) {
    return (
        <UserProvider>
            <BrowserRouter>{children}</BrowserRouter>
        </UserProvider>
    );
}
