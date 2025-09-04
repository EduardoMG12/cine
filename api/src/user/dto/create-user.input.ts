import { InputType, Field, ID } from "@nestjs/graphql";

@InputType()
export class CreateUserInput {
	@Field()
	username: string;

	@Field()
	full_name: string;

	@Field()
	email: string;

	@Field()
	password_hash: string;
}
