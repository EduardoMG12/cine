import { Resolver, Query, Args, ID, Mutation } from "@nestjs/graphql";
import type { UserService } from "./user.service";
import { User } from "./dto/user.model";
import { CreateUserInput } from "./dto/create-user.input";
import { UpdateUserInput } from "./dto/update-user.input";
import { RegisterPayload } from "./dto/register.payload";
import type { AuthService } from "../auth/auth.service";

@Resolver(() => User)
export class UserResolver {
	constructor(
		private readonly userService: UserService,
		private readonly authService: AuthService,
	) {}

	@Query(() => User)
	async findOneUser(@Args("id", { type: () => ID }) id: string): Promise<User> {
		return this.userService.findOne(id);
	}

	@Query(() => [User])
	async findUsers(): Promise<User[]> {
		return this.userService.findAll();
	}

	@Mutation(() => RegisterPayload)
	async register(
		@Args("input", { type: () => CreateUserInput }) input: CreateUserInput,
	): Promise<RegisterPayload> {
		const newUser = await this.userService.create(input);
		const token = await this.authService.generateToken(newUser);
		return { user: newUser, token };
	}

	@Mutation(() => User)
	async updateUser(
		@Args("id", { type: () => ID }) id: string,
		@Args("input", { type: () => UpdateUserInput }) input: UpdateUserInput,
	): Promise<User> {
		return this.userService.update(id, input);
	}

	@Mutation(() => Boolean)
	async removeUser(
		@Args("id", { type: () => ID }) id: string,
	): Promise<boolean> {
		await this.userService.remove(id);
		return true;
	}
}
