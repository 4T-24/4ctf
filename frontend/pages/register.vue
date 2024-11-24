<template>
  <div
    class="border border-solid p-7 border-black-600 mt-24 flex flex-col items-center rounded-lg"
  >
    <h1 class="text-xl text-gray-700 mb-5">Register</h1>

    <Form
      v-slot="$form"
      :resolver
      :validateOnValueUpdate="false"
      :validateOnBlur="true"
      @submit="onFormSubmit"
      class="flex flex-col gap-4 w-full sm:w-80"
    >
      <div class="flex flex-col gap-2">
        <FloatLabel variant="on">
          <IconField
            ><InputIcon class="pi pi-user"></InputIcon>
            <InputText
              id="username"
              name="username"
              type="text"
              fluid
              v-model="userLogin.username"
            />
            <label for="username">Username</label> </IconField
          ><Message
            v-if="$form.username?.invalid"
            severity="error"
            size="small"
            variant="simple"
            >{{ $form.username.error?.message }}</Message
          >
        </FloatLabel>
        <FloatLabel variant="on">
          <IconField>
            <InputIcon class="pi pi-at" />

            <InputText
              id="email"
              name="email"
              type="email"
              fluid
              v-model="userLogin.email"
            />
            <label for="email">Email</label>
          </IconField>

          <Message
            v-if="$form.email?.invalid"
            severity="error"
            size="small"
            variant="simple"
            >{{ $form.email.error?.message }}</Message
          ></FloatLabel
        >
        <FloatLabel variant="on">
          <IconField>
            <InputIcon class="pi pi-lock" />
            <InputText
              id="password"
              name="password"
              type="text"
              fluid
              v-model="userLogin.password"
            />
            <label for="password">Password</label></IconField
          ><Message
            v-if="$form.password?.invalid"
            severity="error"
            size="small"
            variant="simple"
            >{{ $form.password.error?.message }}</Message
          >
        </FloatLabel>
      </div>
      <Button type="submit" severity="secondary" label="Submit" />
    </Form>
    <p class="text-gray-600 mt-4">
      If you already got an account
      <NuxtLink class="text-green-500 underline" to="login">login</NuxtLink>
    </p>
  </div>
  <Toast />
</template>

<script lang="ts" setup>
import FloatLabel from "primevue/floatlabel";
import Button from "primevue/button";
import { useLoginStore } from "~/store/login";

const toast = useToast();
const loginStore = useLoginStore();
const userLogin = ref<{
  email: string | null;
  username: string | null;
  password: string | null;
}>({
  username: null,
  email: null,
  password: null,
});

type FormError = {
  username: Error[];
  email: Error[];
  password: Error[];
};

type Error = {
  message: string;
};

const resolver = ({
  values,
}: {
  values: { username: string; email: string; password: string };
}) => {
  const errors: FormError = { username: [], email: [], password: [] };

  if (!values.username) {
    errors.username = [{ message: "Username is required." }];
  }

  if (!values.email) {
    errors.email = [{ message: "Email is required." }];
  }

  if (!values.password) {
    errors.password = [{ message: "Password is required." }];
  } else {
    if (values.password.length < 12) {
      errors.password = [
        { message: "Password must be at least 12 characters long" },
      ];
    }
  }
  return {
    errors,
  };
};
const onFormSubmit = async (e:{valid:boolean}) => {
  if (e.valid) {
    const err = await loginStore.register(userLogin.value);
    if (err.length) {
      toast.add({
        severity: "error",
        summary: "Error",
        detail: `${err.at(1)?.message}: ${err.at(0)?.message}`,
        life: 3000,
      });
    } else {
      loginStore.setFromCreated(true)
      navigateTo("/login");
    }
  }
};
</script>

<style>
.container {
  padding: 2rem;
  height: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
}
.form-field {
  margin: 10px 0;
}
</style>
