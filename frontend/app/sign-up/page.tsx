"use client";

import Auth from "@/app/components/Auth";
import { authUrl } from "@/utils/network";
import axios, { AxiosError } from "axios";
import { FormEvent, useState } from "react";
import  { toast } from "react-toastify";
import { useRouter } from "next/navigation";
import { errorHandler } from "@/utils/errorHandler";

const Register = () => {
    const [loading, setLoading] = useState(false);
    const router = useRouter();

    const onSubmit = async (
        e: FormEvent<HTMLFormElement>,
        formRef: React.RefObject<HTMLFormElement | null>
    ) => {
        e.preventDefault();
        setLoading(true);

        let arg = {
            email: formRef.current?.email.value,
            password: formRef.current?.password.value,
        }

        const response = await axios
            .post(authUrl.register, arg)
            .catch((e: AxiosError) => errorHandler(e));

        setLoading(false);

        if (response) {
            toast("User created successfully", {
                type: "success",
            })

            router.push("/login");
        }
    }

    return <Auth
        title="Sign Up"
        buttonTitle="Sign Up"
        loading={loading}
        onSubmit={onSubmit}
        accountInfoText={{
            initialText: "Already have an account?",
            actionText: "Login",
            actionLink: "/login"
        }}
    />
}

export default Register;
