import { Module } from "@nestjs/common";
import { UserResolver } from "./user.resolver";
import { UserService } from "./user.service";
import { TypeOrmModule } from "@nestjs/typeorm";
import { User } from "../entities/user.entity";
import { AuthModule } from "../auth/auth.module";

@Module({
	imports: [TypeOrmModule.forFeature([User]), AuthModule],
	providers: [UserResolver, UserService],
	exports: [UserService],
})
export class UserModule {}
