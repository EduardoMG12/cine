import { Injectable, NotFoundException } from "@nestjs/common";
import { InjectRepository } from "@nestjs/typeorm";
import type { Repository } from "typeorm";
import { User } from "../entities/user.entity";
import type { CreateUserInput } from "./dto/create-user.input";
import type { UpdateUserInput } from "./dto/update-user.input";
import * as bcrypt from "bcrypt";

@Injectable()
export class UserService {
	constructor(
		@InjectRepository(User)
		private readonly userRepository: Repository<User>,
	) {}

	async create(input: CreateUserInput): Promise<User> {
		const hashedPassword = await bcrypt.hash(input.password_hash, 10);
		const newUser = this.userRepository.create({
			...input,
			password_hash: hashedPassword,
		});
		return this.userRepository.save(newUser);
	}

	async findAll(): Promise<User[]> {
		return this.userRepository.find();
	}

	async findOne(id: string): Promise<User> {
		const user = await this.userRepository.findOne({ where: { id } });
		if (!user) {
			throw new NotFoundException(`User with ID ${id} not found.`);
		}
		return user;
	}

	async update(id: string, input: UpdateUserInput): Promise<User> {
		await this.userRepository.update(id, input);
		return this.findOne(id);
	}

	async remove(id: string): Promise<void> {
		const result = await this.userRepository.delete(id);
		if (result.affected === 0) {
			throw new NotFoundException(`User with ID ${id} not found.`);
		}
	}
}
