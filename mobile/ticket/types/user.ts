import { ApiResponse } from "./api";

export type AuthResponse = ApiResponse<
    { of_account: User; token: string; }
>;

export type User = {
    account_id: number;
    email: string;
    role: string;
	createdAt: string;
	updatedAt: string;
}