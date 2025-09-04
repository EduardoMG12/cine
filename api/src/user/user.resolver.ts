// src/user/user.resolver.ts

import { Resolver, Query, Args, ID, Mutation } from "@nestjs/graphql";
import type { UserService } from "./user.service";
import { User } from "./dto/user.model";
import { User as UserEntity } from "./user.entity";
import type { CreateUserInput } from "./dto/create-user.input"; // Você precisará criar esta classe

@Resolver(() => User)
export class UserResolver {
	constructor(private readonly userService: UserService) {}

	@Query(() => User)
	async findOneUser(@Args("id", { type: () => ID }) id: string): Promise<User> {
		return this.userService.findOne(id);
	}

	@Query(() => [User])
	async findUsers(): Promise<User[]> {
		return this.userService.findAll();
	}

	@Mutation(() => User)
	async createUser(@Args("input") input: CreateUserInput): Promise<User> {
		// Crie e salve o usuário usando o serviço
		return this.userService.create(input);
	}
}
