import { ApiResponse } from "./api";

export enum UserRole {
	Attendee = "attendee",
	Manager = "manager",
}

export type AuthResponse = ApiResponse<
    { of_account: User; access_token: string; refresh_token: string}
>;

export type User = {
    account_id: number;
    email: string;
    role: UserRole;
	createdAt: string;
	updatedAt: string;
}