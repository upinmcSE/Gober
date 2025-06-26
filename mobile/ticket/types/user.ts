import { ApiResponse } from "./api";

export type AuthResponse = ApiResponse<
    { token: string; refreshToken: string; user: User }
>;

export type User = {
    id: number;
    email: string;
    role: string;
	createdAt: string;
	updatedAt: string;
}