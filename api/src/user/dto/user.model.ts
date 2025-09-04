import { ObjectType, Field, ID } from "@nestjs/graphql";

@ObjectType()
export class User {
	@Field(() => ID)
	id: string;

	@Field()
	username: string;

	@Field()
	full_name: string;

	@Field()
	email: string;

	@Field()
	password_hash: string;
}
