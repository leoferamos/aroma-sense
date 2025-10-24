import React from 'react';
import { Link } from 'react-router-dom';
import LegalPageLayout from '../components/LegalPageLayout';

const Privacy: React.FC = () => {
    return (
        <LegalPageLayout title="Privacy Policy" lastUpdate="October 23, 2025">
            <section className="space-y-4 text-justify">
                            <p>
                                A sua privacidade é <strong className="text-gray-900">muito importante</strong> para nós. Esta Política de Privacidade explica como o Aroma Sense, de responsabilidade de Julio Oliveira e Leonardo Ramos, coleta, utiliza, armazena e protege as informações dos usuários que acessam e utilizam nossos serviços.
                            </p>

                            <p>Ao utilizar o site Aroma Sense, você concorda com as práticas descritas nesta política.</p>

                            <h3 className="text-lg font-semibold text-gray-900 mt-6">1. Informações Coletadas</h3>
                            <p>
                                O Aroma Sense coleta apenas as informações estritamente necessárias para o funcionamento do e-commerce e da personalização da experiência do usuário.
                            </p>
                            <p>As informações coletadas incluem:</p>
                            <ul className="list-disc list-inside text-gray-700 ml-3 space-y-1">
                                <li><strong>E-mail e senha:</strong> fornecidos no momento do cadastro.</li>
                                <li><strong>Cookies de navegação:</strong> usados para melhorar a experiência do usuário e personalizar recomendações.</li>
                            </ul>
                            <p className="mt-2">O Aroma Sense não coleta dados sensíveis, como CPF, endereço, telefone ou informações de pagamento diretamente em sua base — esses dados são processados de forma segura pelos provedores de pagamento.</p>

                            <h3 className="text-lg font-semibold text-gray-900 mt-6">2. Finalidade da Coleta</h3>
                            <p>Os dados são coletados para as seguintes finalidades:</p>
                            <ul className="list-disc list-inside text-gray-700 ml-3 space-y-1">
                                <li>Permitir o acesso à conta do usuário (login e autenticação);</li>
                                <li>Gerenciar compras e histórico de pedidos;</li>
                                <li>Oferecer recomendações personalizadas com base nas preferências do usuário;</li>
                                <li>Melhorar a experiência de navegação e personalização do site;</li>
                                <li>Garantir a segurança e integridade da plataforma.</li>
                            </ul>

                            <h3 className="text-lg font-semibold text-gray-900 mt-6">3. Cookies e Tecnologias de Rastreamento</h3>
                            <p>O Aroma Sense utiliza cookies e tecnologias semelhantes para:</p>
                            <ul className="list-disc list-inside text-gray-700 ml-3 space-y-1">
                                <li>Armazenar preferências do usuário;</li>
                                <li>Manter sessões ativas;</li>
                                <li>Gerar estatísticas anônimas de uso do site;</li>
                                <li>Aperfeiçoar as recomendações feitas pela inteligência artificial.</li>
                            </ul>
                            <p className="mt-2">O usuário pode desativar os cookies nas configurações do navegador, mas isso pode limitar certas funcionalidades do site.</p>

                            <h3 className="text-lg font-semibold text-gray-900 mt-6">4. Compartilhamento de Dados</h3>
                            <p>O Aroma Sense não compartilha informações pessoais dos usuários com terceiros, exceto quando necessário para:</p>
                            <ul className="list-disc list-inside text-gray-700 ml-3 space-y-1">
                                <li>Processamento de pagamentos (via parceiros como gateways de pagamento);</li>
                                <li>Entrega de produtos, quando aplicável;</li>
                                <li>Cumprimento de obrigações legais ou regulatórias.</li>
                            </ul>
                            <p className="mt-2">Em nenhum caso os dados são vendidos, alugados ou cedidos para fins de marketing externo.</p>

                            <h3 className="text-lg font-semibold text-gray-900 mt-6">5. Armazenamento e Segurança</h3>
                            <p>
                                Os dados são armazenados em ambientes seguros, com medidas técnicas e organizacionais adequadas para evitar acesso não autorizado, destruição, perda ou alteração indevida.
                            </p>
                            <p className="mt-2">As senhas são criptografadas e não são visualizadas por nossa equipe. Apesar dos esforços, nenhum sistema é 100% seguro. Em caso de incidente de segurança que comprometa dados pessoais, o Aroma Sense notificará os usuários afetados conforme exigido pela LGPD.</p>

                            <h3 className="text-lg font-semibold text-gray-900 mt-6">6. Direitos do Usuário (LGPD)</h3>
                            <p>De acordo com a Lei Geral de Proteção de Dados (LGPD), o usuário tem direito a:</p>
                            <ul className="list-disc list-inside text-gray-700 ml-3 space-y-1">
                                <li>Acessar os dados armazenados sobre si;</li>
                                <li>Corrigir dados incompletos ou incorretos;</li>
                                <li>Solicitar a exclusão da conta e dos dados pessoais;</li>
                                <li>Revogar o consentimento para uso dos dados;</li>
                                <li>Solicitar informações sobre o compartilhamento de dados.</li>
                            </ul>
                            <p className="mt-2">Para exercer qualquer um desses direitos, basta entrar em contato pelo e-mail:
                                <a href="mailto:suporte.aromasene@gmail.com" className="text-blue-600 hover:underline ml-1">📩 suporte.aromasene@gmail.com</a>
                            </p>

                            <h3 className="text-lg font-semibold text-gray-900 mt-6">7. Retenção e Exclusão de Dados</h3>
                            <p>Os dados do usuário serão mantidos apenas pelo tempo necessário para:</p>
                            <ul className="list-disc list-inside text-gray-700 ml-3 space-y-1">
                                <li>Cumprir as finalidades descritas nesta Política;</li>
                                <li>Atender exigências legais ou contratuais.</li>
                            </ul>
                            <p className="mt-2">Após o encerramento da conta ou solicitação de exclusão, os dados serão removidos de forma segura e definitiva.</p>

                            <h3 className="text-lg font-semibold text-gray-900 mt-6">8. Inteligência Artificial e Personalização</h3>
                            <p>O Aroma Sense utiliza algoritmos de inteligência artificial para sugerir perfumes com base nas preferências e interações do usuário. Essas recomendações são automáticas e não envolvem decisões humanas diretas. Nenhuma decisão de caráter legal, financeiro ou pessoal é tomada exclusivamente por meio da IA.</p>

                            <h3 className="text-lg font-semibold text-gray-900 mt-6">9. Alterações na Política</h3>
                            <p>Esta Política de Privacidade pode ser atualizada periodicamente para refletir melhorias ou mudanças legais. Recomenda-se a leitura regular desta página para se manter informado sobre eventuais alterações.</p>

                            <h3 className="text-lg font-semibold text-gray-900 mt-6">10. Contato</h3>
                            <p>Em caso de dúvidas, solicitações ou reclamações sobre o tratamento de dados, entre em contato com nossa equipe pelo e-mail:
                                <a href="mailto:suporte.aromasene@gmail.com" className="text-blue-600 hover:underline ml-1">📩 suporte.aromasene@gmail.com</a>
                            </p>

                            <p className="text-center text-gray-500 text-sm mt-6">Aroma Sense © 2025 — Todos os direitos reservados.</p>
                        </section>

                        <div className="mt-6">
                            <Link
                                to="/register"
                                className="block w-full md:w-1/2 mx-auto text-white text-lg font-medium py-3 rounded-full mt-2 bg-blue-600 hover:bg-blue-700 text-center transition-colors"
                            >
                                Voltar
                            </Link>
                        </div>
        </LegalPageLayout>
    );
};

export default Privacy;
