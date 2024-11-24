import { defineStore } from "pinia";
import config from "@/data/config.yaml";

export const useLoginStore = defineStore({
  id: "LoginStore",
  state: () => ({ isLogged: false, fromCreated: false }),
  actions: {
    setIsLogged(value: boolean) {
      this.isLogged = value;
    },
    setFromCreated(value: boolean){
      this.fromCreated = value
    },
    async login(data: { username: string|null; password: string|null }) {
      const resp = await fetch(`${config.serverURL}/auth/login`, {
        method: "POST",
        body: JSON.stringify(data),
      });
      console.debug(resp)
      const json = await resp.json();
      if (json.errors) {
        return json.errors as Error[];
      } else {
        return [];
      }
    },
    async register(data: {
      username: string | null;
      email: string | null;
      password: string | null;
    }) {
      const resp = await fetch(`${config.serverURL}/auth/register`, {
        method: "POST",
        body: JSON.stringify(data),
      });
      const json = await resp.json();
      if (json.errors) {
        return json.errors as Error[];
      } else {
        return [];
      }
    },
  },
});

type Error = {
  field: string;
  message: string;
};
