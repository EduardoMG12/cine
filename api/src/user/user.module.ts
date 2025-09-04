// src/user/user.module.ts

import { Module } from "@nestjs/common";
import { UserResolver } from "./user.resolver";
import { UserService } from "./user.service";
import { TypeOrmModule } from "@nestjs/typeorm";
import { User } from "../entities/user.entity";

@Module({
	imports: [TypeOrmModule.forFeature([User])], // Importe as entidades aqui
	providers: [UserResolver, UserService],
	exports: [UserService], // Exporte o serviço se ele for usado em outros módulos
})
export class UserModule {}
