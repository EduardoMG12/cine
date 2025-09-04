import { Field, ObjectType } from "@nestjs/graphql";
import { User } from "./user.model";

@ObjectType()
export class RegisterPayload {
	@Field(() => User)
	user: User;

	@Field()
	token: string;
}
