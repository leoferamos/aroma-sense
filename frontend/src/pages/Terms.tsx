import React from 'react';
import LegalPageLayout from '../components/LegalPageLayout';
import AuthBackLink from '../components/AuthBackLink';
import { useTranslation } from 'react-i18next';

const Terms: React.FC = () => {
    const { t } = useTranslation('legal');
    return (
        <LegalPageLayout title={t('title')} lastUpdate="November 30, 2025">
            <h1 className="text-center text-2xl font-semibold text-gray-900 mb-2">{t('heading')}</h1>
            <p className="text-center text-sm text-gray-500 mb-8">{t('lastUpdate')}</p>

            <section className="space-y-4 text-justify">
                            <p>
                                <strong className="text-gray-900">{t('welcome')}</strong>
                            </p>

                            <p>
                                {t('agreement')}
                            </p>

                            <h3 className="text-lg font-semibold text-gray-900 mt-6">{t('section1.title')}</h3>
                            <p>
                                {t('section1.content')}
                            </p>

                            <h3 className="text-lg font-semibold text-gray-900 mt-6">{t('section2.title')}</h3>
                            <p>
                                {t('section2.content')}
                            </p>
                            <ul className="list-disc list-inside text-gray-700 ml-3 space-y-1">
                                {t('section2.responsibilities', { returnObjects: true }).map((item: string, index: number) => (
                                    <li key={index}>{item}</li>
                                ))}
                            </ul>
                            <p>
                                {t('section2.liability')}
                            </p>

                            <h3 className="text-lg font-semibold text-gray-900 mt-6">{t('section3.title')}</h3>
                            <p>{t('section3.intro')}</p>
                            <ul className="list-disc list-inside text-gray-700 ml-3 space-y-1">
                                {t('section3.features', { returnObjects: true }).map((item: string, index: number) => (
                                    <li key={index}>{item}</li>
                                ))}
                            </ul>
                            <p>
                                {t('section3.modifications')}
                            </p>

                            <h3 className="text-lg font-semibold text-gray-900 mt-6">{t('section4.title')}</h3>
                            <p>
                                {t('section4.content')}
                            </p>
                            <p>{t('section4.notification')}</p>

                            <h3 className="text-lg font-semibold text-gray-900 mt-6">{t('section5.title')}</h3>
                            <p>
                                {t('section5.content')}
                                <a href={`mailto:${t('section5.email')}`} className="text-blue-600 hover:underline ml-1">
                                    {t('section5.email')}
                                </a>.
                            </p>

                            <h3 className="text-lg font-semibold text-gray-900 mt-6">{t('section6.title')}</h3>
                            <p>
                                {t('section6.content')}
                            </p>

                            <h3 className="text-lg font-semibold text-gray-900 mt-6">{t('section7.title')}</h3>
                            <p>
                                {t('section7.content')}
                            </p>
                            <p>{t('section7.acknowledgment')}</p>
                            <ul className="list-disc list-inside text-gray-700 ml-3 space-y-1">
                                {t('section7.points', { returnObjects: true }).map((item: string, index: number) => (
                                    <li key={index}>{item}</li>
                                ))}
                            </ul>

                            <h3 className="text-lg font-semibold text-gray-900 mt-6">{t('section8.title')}</h3>
                            <p>
                                {t('section8.content')}
                            </p>

                            <h3 className="text-lg font-semibold text-gray-900 mt-6">{t('section9.title')}</h3>
                            <p>
                                {t('section9.content')}
                            </p>
                            <ul className="list-disc list-inside text-gray-700 ml-3 space-y-1">
                                {t('section9.situations', { returnObjects: true }).map((item: string, index: number) => (
                                    <li key={index}>{item}</li>
                                ))}
                            </ul>
                            <p className="mt-3"><strong>{t('section9.process')}</strong></p>
                            <ol className="list-decimal list-inside text-gray-700 ml-3 space-y-1">
                                {t('section9.processSteps', { returnObjects: true }).map((item: string, index: number) => (
                                    <li key={index}>{item}</li>
                                ))}
                            </ol>
                            <p className="mt-3"><strong>{t('section9.consequences')}</strong></p>
                            <ul className="list-disc list-inside text-gray-700 ml-3 space-y-1">
                                {t('section9.consequencesList', { returnObjects: true }).map((item: string, index: number) => (
                                    <li key={index}>{item}</li>
                                ))}
                            </ul>
                            <p className="mt-2">
                                {t('section9.commitment')}
                            </p>

                            <h3 className="text-lg font-semibold text-gray-900 mt-6">{t('section10.title')}</h3>
                            <p>
                                {t('section10.content')}
                            </p>

                            <h3 className="text-lg font-semibold text-gray-900 mt-6">{t('section11.title')}</h3>
                            <p>
                                {t('section11.content')}
                                <a href={`mailto:${t('section11.email')}`} className="text-blue-600 hover:underline ml-1">
                                    {t('section11.email')}
                                </a>.
                            </p>

                            <p className="text-center text-gray-500 text-sm mt-6">
                                {t('copyright')}
                            </p>
                        </section>

                        <div className="mt-6">
                            <AuthBackLink />
                        </div>
        </LegalPageLayout>
    );
};

export default Terms;
