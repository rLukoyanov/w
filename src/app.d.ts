declare global {
	namespace App {
		interface Locals {
			user: import('pocketbase').RecordModel | null;
		}
	}
}

export {};
