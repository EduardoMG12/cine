import { Injectable } from "@nestjs/common";
import type { JwtService } from "@nestjs/jwt";
import type { User } from "../entities/user.entity";

@Injectable()
export class AuthService {
	constructor(private readonly jwtService: JwtService) {}

	async generateToken(user: User): Promise<string> {
		const payload = { username: user.username, sub: user.id };
		return this.jwtService.sign(payload);
	}
}
