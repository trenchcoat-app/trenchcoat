import { useTranslation } from "react-i18next";
import { Link } from "@tanstack/react-router";

import notFoundStyles from "./NotFound.module.css";
import pageStyles from "./pages.module.css";

export const NotFound = () => {
    const { t } = useTranslation();

    return (
        <div className={`${pageStyles.pageWrapper} ${pageStyles.pageBottomPadding} ${pageStyles.pageCenterContent}`}>
            <div className={notFoundStyles.wrapper}>
                <p className={notFoundStyles.title}>404</p>
                <p className={notFoundStyles.message}>{t("PAGE_NOT_FOUND_MESSAGE")}</p>
                <Link to="/">{t("BACK_TO_HOME")}</Link>
            </div>
        </div>
    );
};
