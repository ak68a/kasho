"use client";

import Auth from "@/app/components/Auth";
import axios from "axios";
import { errorHandler } from "@/utils/errorHandler";
import { AxiosError } from "axios";
import { FormEvent, useState } from "react";
import  { toast } from "react-toastify";
import { authUrl } from "@/utils/network";
import { useRouter } from "next/navigation";

const Login = () => {
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
            .post(authUrl.login, arg)
            .catch((e: AxiosError) => errorHandler(e));

        setLoading(false);

        if (response) {
            toast("You're logged in!", {
                type: "success",
            })

            router.push("/");
        }
    };

    return <Auth showRemembered loading={loading} onSubmit={onSubmit} />

};

export default Login;
