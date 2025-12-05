import React from 'react';
import LegalPageLayout from '../components/LegalPageLayout';
import AuthBackLink from '../components/AuthBackLink';
import { useTranslation } from 'react-i18next';

const Privacy: React.FC = () => {
    const { t } = useTranslation('privacy');
    return (
        <LegalPageLayout title={t('title')} lastUpdate="November 30, 2025">
            <section className="space-y-4 text-justify">
                            <p>
                                {t('intro')}
                            </p>

                            <p>{t('agreement')}</p>

                            <h3 className="text-lg font-semibold text-gray-900 mt-6">{t('section1.title')}</h3>
                            <p>
                                {t('section1.content')}
                            </p>
                            <p>{t('section1.collected')}</p>
                            <ul className="list-disc list-inside text-gray-700 ml-3 space-y-1">
                                {t('section1.items', { returnObjects: true }).map((item: string, index: number) => (
                                    <li key={index} dangerouslySetInnerHTML={{ __html: item }}></li>
                                ))}
                            </ul>
                            <p className="mt-2">{t('section1.sensitive')}</p>

                            <h3 className="text-lg font-semibold text-gray-900 mt-6">{t('section2.title')}</h3>
                            <p>{t('section2.intro')}</p>
                            <ul className="list-disc list-inside text-gray-700 ml-3 space-y-1">
                                {t('section2.purposes', { returnObjects: true }).map((item: string, index: number) => (
                                    <li key={index}>{item}</li>
                                ))}
                            </ul>

                            <h3 className="text-lg font-semibold text-gray-900 mt-6">{t('section3.title')}</h3>
                            <p>{t('section3.intro')}</p>
                            <ul className="list-disc list-inside text-gray-700 ml-3 space-y-1">
                                {t('section3.uses', { returnObjects: true }).map((item: string, index: number) => (
                                    <li key={index}>{item}</li>
                                ))}
                            </ul>
                            <p className="mt-2">{t('section3.control')}</p>

                            <h3 className="text-lg font-semibold text-gray-900 mt-6">{t('section4.title')}</h3>
                            <p>{t('section4.content')}</p>
                            <ul className="list-disc list-inside text-gray-700 ml-3 space-y-1">
                                {t('section4.exceptions', { returnObjects: true }).map((item: string, index: number) => (
                                    <li key={index}>{item}</li>
                                ))}
                            </ul>
                            <p className="mt-2">{t('section4.noSale')}</p>

                            <h3 className="text-lg font-semibold text-gray-900 mt-6">{t('section5.title')}</h3>
                            <p>
                                {t('section5.content')}
                            </p>
                            <p className="mt-2">{t('section5.passwords')}</p>

                            <h3 className="text-lg font-semibold text-gray-900 mt-6">{t('section6.title')}</h3>
                            <p>{t('section6.intro')}</p>
                            <ul className="list-disc list-inside text-gray-700 ml-3 space-y-1">
                                {t('section6.rights', { returnObjects: true }).map((item: string, index: number) => (
                                    <li key={index}>{item}</li>
                                ))}
                            </ul>
                            <p className="mt-2"><strong>{t('section6.deactivationRights.title')}</strong></p>
                            <p>
                                {t('section6.deactivationRights.content')}
                            </p>
                            <ul className="list-disc list-inside text-gray-700 ml-3 space-y-1">
                                {t('section6.deactivationRights.rights', { returnObjects: true }).map((item: string, index: number) => (
                                    <li key={index}>{item}</li>
                                ))}
                            </ul>
                            <p className="mt-2">{t('section6.contact')}</p>

                            <h3 className="text-lg font-semibold text-gray-900 mt-6">{t('section7.title')}</h3>
                            <p>{t('section7.intro')}</p>
                            <ul className="list-disc list-inside text-gray-700 ml-3 space-y-1">
                                {t('section7.purposes', { returnObjects: true }).map((item: string, index: number) => (
                                    <li key={index}>{item}</li>
                                ))}
                            </ul>
                            <p className="mt-2"><strong>{t('section7.platformDeactivation.title')}</strong></p>
                            <p>
                                {t('section7.platformDeactivation.content')}
                            </p>
                            <p className="mt-2">{t('section7.postClosure')}</p>

                            <h3 className="text-lg font-semibold text-gray-900 mt-6">{t('section8.title')}</h3>
                            <p>{t('section8.content')}</p>

                            <h3 className="text-lg font-semibold text-gray-900 mt-6">{t('section9.title')}</h3>
                            <p>{t('section9.content')}</p>

                            <h3 className="text-lg font-semibold text-gray-900 mt-6">{t('section10.title')}</h3>
                            <p>{t('section10.content')}</p>

                            <p className="text-center text-gray-500 text-sm mt-6">{t('copyright')}</p>
                        </section>

                        <div className="mt-6">
                            <AuthBackLink />
                        </div>
        </LegalPageLayout>
    );
};

export default Privacy;
