import { Suspense, lazy } from "react";
import type { ClassKey } from "keycloakify/login";
import type { KcContext } from "./KcContext";
import { useI18n } from "./i18n";
import DefaultPage from "keycloakify/login/DefaultPage";
import Template from "keycloakify/login/Template";
import { getKcClsx } from "keycloakify/login/lib/kcClsx";
import "./KcPage.css";

const UserProfileFormFields = lazy(
    () => import("keycloakify/login/UserProfileFormFields")
);

const doMakeUserConfirmPassword = true;

function CustomTemplate(props: Parameters<typeof Template>[0]) {
    const { kcContext, doUseDefaultCss, classes, children } = props;
    const { kcClsx } = getKcClsx({ doUseDefaultCss, classes });
    const { url, realm, messagesPerField, social } = kcContext;
    console.error("social:", social.providers);

    return (
        <div
            style={{
                position: "fixed",
                top: 0,
                left: 0,
                right: 0,
                bottom: 0,
                overflow: "auto"
            }}
        >
            {/* Blurred background */}
            <div
                style={{
                    position: "fixed",
                    top: 0,
                    left: 0,
                    right: 0,
                    bottom: 0,
                    backgroundImage: "url('./background.png')",
                    backgroundSize: "cover",
                    backgroundPosition: "center",
                    filter: "blur(6px)",
                    zIndex: -1
                }}
            />

            {/* Content wrapper */}
            <div
                style={{
                    minHeight: "100vh",
                    display: "flex",
                    alignItems: "center",
                    justifyContent: "center",
                    padding: "20px"
                }}
            >
                {/* Login form container */}
                <div
                    style={{
                        background: "rgba(30, 30, 30, 0.95)",
                        borderRadius: "12px",
                        padding: "40px",
                        maxWidth: "400px",
                        width: "100%",
                        boxShadow: "0 8px 32px rgba(0, 0, 0, 0.5)",
                        position: "relative"
                    }}
                >
                    <h1 className={kcClsx("kcHeaderClass")}>Sign In</h1>

                    {kcContext.pageId === "login.ftl" && (
                        <>
                            {(realm as any).password && social?.providers && (
                                <>
                                    <div
                                        className={kcClsx(
                                            "kcFormSocialAccountSectionClass"
                                        )}
                                    >
                                        <ul
                                            className={kcClsx(
                                                "kcFormSocialAccountListClass"
                                            )}
                                        >
                                            {social.providers.map(provider => {
                                                const getProviderLogo = (
                                                    providerId: string
                                                ) => {
                                                    const logos: Record<string, string> =
                                                        {
                                                            google: "./google-logo.png",
                                                            intra42: "./42-logo.png"
                                                        };
                                                    return logos[providerId] || null;
                                                };

                                                return (
                                                    <li key={provider.providerId}>
                                                        <a
                                                            href={provider.loginUrl}
                                                            className={kcClsx(
                                                                "kcFormSocialAccountListButtonClass"
                                                            )}
                                                        >
                                                            <span
                                                                className={kcClsx(
                                                                    "kcFormSocialAccountNameClass"
                                                                )}
                                                            >
                                                                Sign in with{" "}
                                                                {getProviderLogo(
                                                                    provider.displayName
                                                                ) ? (
                                                                    <img
                                                                        src={
                                                                            getProviderLogo(
                                                                                provider.displayName
                                                                            )!
                                                                        }
                                                                        alt={
                                                                            provider.displayName
                                                                        }
                                                                        className="kc-custom-social-logo"
                                                                    />
                                                                ) : (
                                                                    provider.displayName
                                                                )}
                                                            </span>
                                                        </a>
                                                    </li>
                                                );
                                            })}
                                        </ul>
                                    </div>
                                    <div className="kc-custom-divider">or</div>
                                </>
                            )}

                            <form action={url.loginAction} method="post">
                                <div className={kcClsx("kcFormGroupClass")}>
                                    <label
                                        htmlFor="username"
                                        className={kcClsx("kcLabelClass")}
                                    >
                                        Username or Email
                                    </label>
                                    <input
                                        type="text"
                                        id="username"
                                        name="username"
                                        className={kcClsx("kcInputClass")}
                                        autoFocus
                                        autoComplete="username"
                                    />
                                </div>

                                <div className={kcClsx("kcFormGroupClass")}>
                                    <label
                                        htmlFor="password"
                                        className={kcClsx("kcLabelClass")}
                                    >
                                        Password
                                    </label>
                                    <input
                                        type="password"
                                        id="password"
                                        name="password"
                                        className={kcClsx("kcInputClass")}
                                        autoComplete="current-password"
                                    />
                                </div>

                                <div className={kcClsx("kcFormGroupClass")}>
                                    <div className={kcClsx("kcFormSettingClass")}>
                                        <div className={kcClsx("kcFormOptionsClass")}>
                                            {(realm as any).rememberMe && (
                                                <>
                                                    <input
                                                        type="checkbox"
                                                        id="rememberMe"
                                                        name="rememberMe"
                                                    />
                                                    <label htmlFor="rememberMe">
                                                        Remember me
                                                    </label>
                                                </>
                                            )}
                                        </div>
                                        {realm.resetPasswordAllowed && (
                                            <div
                                                className={kcClsx(
                                                    "kcFormOptionsWrapperClass"
                                                )}
                                            >
                                                <a href={url.loginResetCredentialsUrl}>
                                                    Forgot Password?
                                                </a>
                                            </div>
                                        )}
                                    </div>
                                </div>

                                <div className={kcClsx("kcFormButtonsClass")}>
                                    <button
                                        type="submit"
                                        className={kcClsx(
                                            "kcButtonClass",
                                            "kcButtonPrimaryClass",
                                            "kcButtonBlockClass",
                                            "kcButtonLargeClass"
                                        )}
                                    >
                                        Sign In
                                    </button>
                                </div>
                            </form>
                        </>
                    )}

                    {kcContext.pageId !== "login.ftl" && children}
                </div>
            </div>
        </div>
    );
}

export default function KcPage(props: { kcContext: KcContext }) {
    const { kcContext } = props;

    const { i18n } = useI18n({ kcContext });

    return (
        <Suspense>
            {(() => {
                switch (kcContext.pageId) {
                    default:
                        return (
                            <DefaultPage
                                kcContext={kcContext}
                                i18n={i18n}
                                classes={classes}
                                Template={CustomTemplate}
                                doUseDefaultCss={false}
                                UserProfileFormFields={UserProfileFormFields}
                                doMakeUserConfirmPassword={doMakeUserConfirmPassword}
                            />
                        );
                }
            })()}
        </Suspense>
    );
}

const classes = {
    kcFormGroupClass: "kc-custom-form-group",
    kcLabelClass: "kc-custom-label",
    kcInputClass: "kc-custom-input",
    kcFormButtonsClass: "kc-custom-form-buttons",
    kcButtonClass: "kc-custom-button",
    kcButtonPrimaryClass: "kc-custom-button-primary",
    kcButtonBlockClass: "kc-custom-button-block",
    kcButtonLargeClass: "kc-custom-button-large",
    kcInputGroup: "kc-custom-input-group",
    kcFormPasswordVisibilityButtonClass: "kc-custom-password-visibility",
    kcFormSettingClass: "kc-custom-form-options",
    kcFormOptionsClass: "kc-custom-form-options-left",
    kcFormOptionsWrapperClass: "kc-custom-form-options-right",
    kcHeaderClass: "kc-custom-header",
    kcAlertClass: "kc-custom-alert",
    kcAlertErrorClass: "kc-custom-alert-error",
    kcAlertWarningClass: "kc-custom-alert-warning",
    kcAlertSuccessClass: "kc-custom-alert-success",
    kcAlertInfoClass: "kc-custom-alert-info",
    kcFormSocialAccountSectionClass: "kc-custom-social-section",
    kcFormSocialAccountListClass: "kc-custom-social-list",
    kcFormSocialAccountListButtonClass: "kc-custom-social-button",
    kcFormSocialAccountNameClass: "kc-custom-social-name",
    kcFormSocialAccountLinkClass: "kc-custom-social-link"
} satisfies Partial<{ [key in ClassKey]: string }>;
