import React from 'react';
import LegalPageLayout from '../components/LegalPageLayout';
import AuthBackLink from '../components/AuthBackLink';

const Terms: React.FC = () => {
    return (
        <LegalPageLayout title="Terms of Service" lastUpdate="November 27, 2025">
            <h1 className="text-center text-2xl font-semibold text-gray-900 mb-2">Termos de Uso — Aroma Sense</h1>
            <p className="text-center text-sm text-gray-500 mb-8">Última atualização: 30 de novembro de 2025</p>

            <section className="space-y-4 text-justify">
                            <p>
                                <strong className="text-gray-900">Bem-vindo ao Aroma Sense!</strong> Estes Termos de Uso (“Termos”)
                                regulam o acesso e o uso do site e da plataforma Aroma Sense, de propriedade e responsabilidade de
                                Julio Oliveira e Leonardo Ramos (“nós”, “nossa equipe” ou “Aroma Sense”).
                            </p>

                            <p>
                                Ao acessar ou utilizar o site, você concorda integralmente com os presentes Termos. Caso não concorde,
                                recomendamos que não utilize o serviço.
                            </p>

                            <h3 className="text-lg font-semibold text-gray-900 mt-6">1. Objetivo da Plataforma</h3>
                            <p>
                                O Aroma Sense é uma plataforma de e-commerce especializada em perfumes, que oferece ao usuário a
                                possibilidade de explorar, avaliar e adquirir produtos diretamente pelo site, com pagamento via Pix ou
                                cartão de crédito. A plataforma também oferece um sistema de recomendação personalizado com
                                inteligência artificial, sugerindo perfumes de acordo com as preferências do usuário.
                            </p>

                            <h3 className="text-lg font-semibold text-gray-900 mt-6">2. Cadastro e Conta de Usuário</h3>
                            <p>
                                Para realizar compras e utilizar todos os recursos da plataforma, o usuário deverá criar uma conta
                                informando e-mail e senha válidos. O usuário é responsável por:
                            </p>
                            <ul className="list-disc list-inside text-gray-700 ml-3 space-y-1">
                                <li>Manter a confidencialidade de suas credenciais;</li>
                                <li>Fornecer informações verdadeiras e atualizadas;</li>
                                <li>Notificar o Aroma Sense caso suspeite de uso indevido de sua conta.</li>
                            </ul>
                            <p>
                                O Aroma Sense não se responsabiliza por acessos não autorizados decorrentes de negligência na
                                proteção das credenciais.
                            </p>

                            <h3 className="text-lg font-semibold text-gray-900 mt-6">3. Funcionalidades da Plataforma</h3>
                            <p>O usuário poderá:</p>
                            <ul className="list-disc list-inside text-gray-700 ml-3 space-y-1">
                                <li>Navegar e visualizar produtos disponíveis;</li>
                                <li>Realizar compras com pagamento via Pix ou cartão;</li>
                                <li>Avaliar produtos e deixar comentários;</li>
                                <li>Receber recomendações personalizadas por meio de IA;</li>
                                <li>Acompanhar o status de suas compras e entregas.</li>
                            </ul>
                            <p>
                                O Aroma Sense reserva-se o direito de modificar, suspender ou encerrar qualquer funcionalidade a
                                qualquer momento, mediante aviso prévio quando aplicável.
                            </p>

                            <h3 className="text-lg font-semibold text-gray-900 mt-6">4. Compras, Pagamentos e Entregas</h3>
                            <p>
                                As compras realizadas no Aroma Sense estão sujeitas à disponibilidade de estoque. O envio dos
                                produtos é feito em até 30 (trinta) dias corridos, contados a partir da confirmação do pagamento. As
                                formas de pagamento aceitas são Pix e cartão de crédito, processadas de maneira segura por nossos
                                parceiros de pagamento.
                            </p>
                            <p>O usuário receberá informações sobre o andamento do pedido por e-mail.</p>

                            <h3 className="text-lg font-semibold text-gray-900 mt-6">5. Devoluções e Reembolsos</h3>
                            <p>
                                Conforme o Código de Defesa do Consumidor, o usuário poderá solicitar a devolução ou reembolso do
                                produto em até 7 (sete) dias após o recebimento, desde que o produto esteja lacrado e em perfeitas
                                condições. O procedimento de devolução será informado pela equipe de suporte através do e-mail
                                <a href="mailto:suporte.aromasene@gmail.com" className="text-blue-600 hover:underline ml-1">
                                    suporte.aromasene@gmail.com
                                </a>.
                            </p>

                            <h3 className="text-lg font-semibold text-gray-900 mt-6">6. Propriedade Intelectual</h3>
                            <p>
                                Todo o conteúdo do Aroma Sense — incluindo textos, imagens, descrições, logotipos, marcas e design
                                da plataforma — é protegido por direitos autorais e não pode ser copiado, reproduzido, modificado ou
                                distribuído sem autorização prévia e expressa. É proibido o uso indevido de qualquer material
                                disponibilizado no site.
                            </p>

                            <h3 className="text-lg font-semibold text-gray-900 mt-6">7. Limitação de Responsabilidade</h3>
                            <p>
                                O Aroma Sense se esforça para manter o site seguro, estável e atualizado, mas não garante que o
                                funcionamento será ininterrupto ou livre de erros.
                            </p>
                            <p>O usuário reconhece que:</p>
                            <ul className="list-disc list-inside text-gray-700 ml-3 space-y-1">
                                <li>As recomendações feitas por IA têm caráter informativo e não substituem opinião profissional;</li>
                                <li>
                                    O Aroma Sense não se responsabiliza por eventuais falhas de sistema, indisponibilidades temporárias
                                    ou perdas de dados;
                                </li>
                                <li>O uso do site é feito sob total responsabilidade do usuário.</li>
                            </ul>

                            <h3 className="text-lg font-semibold text-gray-900 mt-6">8. Cookies e Tecnologias de Rastreamento</h3>
                            <p>
                                O Aroma Sense utiliza cookies para melhorar a experiência de navegação, armazenar preferências e
                                coletar dados de acesso anônimos. Ao utilizar o site, o usuário concorda com o uso de cookies conforme
                                descrito em nossa Política de Privacidade.
                            </p>

                            <h3 className="text-lg font-semibold text-gray-900 mt-6">9. Suspensão ou Encerramento de Conta</h3>
                            <p>
                                O Aroma Sense poderá suspender temporariamente ou excluir permanentemente contas de usuários nas seguintes situações:
                            </p>
                            <ul className="list-disc list-inside text-gray-700 ml-3 space-y-1">
                                <li><strong>Violação dos Termos de Uso:</strong> Uso indevido da plataforma, spam, conteúdo ofensivo ou discriminatório;</li>
                                <li><strong>Fraude ou Tentativa de Fraude:</strong> Uso de dados falsos, tentativas de burlar pagamentos ou comprometer transações;</li>
                                <li><strong>Conteúdo Inadequado:</strong> Publicação de comentários abusivos, discriminatórios, ilegais ou que violem direitos de terceiros;</li>
                                <li><strong>Comprometimento da Segurança:</strong> Tentativas de hacking, exploração de vulnerabilidades ou uso de ferramentas automatizadas sem autorização;</li>
                                <li><strong>Inatividade Prolongada:</strong> Contas abandonadas sem uso por período superior a 2 anos;</li>
                                <li><strong>Solicitação Judicial:</strong> Cumprimento de ordens legais, decisões judiciais ou requisições de autoridades competentes.</li>
                            </ul>
                            <p className="mt-3"><strong>Processo de Desativação:</strong></p>
                            <ol className="list-decimal list-inside text-gray-700 ml-3 space-y-1">
                                <li>Notificação prévia por e-mail com justificativa detalhada da medida;</li>
                                <li>Período de 7 dias corridos para contestação e apresentação de defesa;</li>
                                <li>Análise da contestação pela equipe de suporte em até 5 dias úteis;</li>
                                <li>Possibilidade de recurso à decisão junto ao suporte (<a href="mailto:suporte.aromasense@gmail.com" className="text-blue-600 hover:underline">suporte.aromasense@gmail.com</a>);</li>
                                <li>Em casos de violação grave ou risco iminente, desativação imediata pode ocorrer sem prévia notificação.</li>
                            </ol>
                            <p className="mt-3"><strong>Consequências:</strong></p>
                            <ul className="list-disc list-inside text-gray-700 ml-3 space-y-1">
                                <li><strong>Suspensão temporária:</strong> Acesso bloqueado por período determinado, com possibilidade de reativação após correção das irregularidades;</li>
                                <li><strong>Exclusão permanente:</strong> Remoção definitiva da conta e dados pessoais após período de retenção de 2 anos (exceto dados necessários para compliance legal, que serão mantidos por até 5 anos). O fluxo de exclusão segue o processo de pedido de exclusão, período de reflexão de 7 dias, confirmação automática pelo sistema (se não houver cancelamento) e posterior anonimização dos dados após o prazo de retenção.</li>
                            </ul>
                            <p className="mt-2">
                                O Aroma Sense compromete-se a agir de forma proporcional, transparente e fundamentada em todos os casos de desativação de contas.
                            </p>

                            <h3 className="text-lg font-semibold text-gray-900 mt-6">10. Alterações nos Termos</h3>
                            <p>
                                O Aroma Sense poderá alterar estes Termos a qualquer momento, mediante atualização nesta página com
                                nova data de vigência. O uso contínuo da plataforma após as alterações constitui aceitação dos novos
                                Termos.
                            </p>

                            <h3 className="text-lg font-semibold text-gray-900 mt-6">11. Contato</h3>
                            <p>
                                Em caso de dúvidas, sugestões ou solicitações, entre em contato pelo e-mail:
                                <a href="mailto:suporte.aromasense@gmail.com" className="text-blue-600 hover:underline ml-1">
                                    suporte.aromasense@gmail.com
                                </a>.
                            </p>

                            <p className="text-center text-gray-500 text-sm mt-6">
                                Aroma Sense © 2025 — Todos os direitos reservados.
                            </p>
                        </section>

                        <div className="mt-6">
                            <AuthBackLink />
                        </div>
        </LegalPageLayout>
    );
};

export default Terms;
