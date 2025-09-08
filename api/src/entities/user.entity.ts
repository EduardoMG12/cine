import { Entity, PrimaryGeneratedColumn, Column } from "typeorm";

@Entity()
export class User {
	@PrimaryGeneratedColumn("uuid")
	id: string;

	@Column({ unique: true })
	username: string;

	@Column()
	full_name: string;

	@Column({ unique: true })
	email: string;

	@Column()
	password_hash: string;
}
